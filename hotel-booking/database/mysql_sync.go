package database

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"hotel-booking/models"
)

type MySQLSync struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewMySQLSyncFromEnv() *MySQLSync {
	if strings.ToLower(os.Getenv("MYSQL_SYNC_ENABLED")) != "true" {
		return nil
	}
	s := &MySQLSync{
		Host:     getenv("MYSQL_HOST", "127.0.0.1"),
		Port:     getenv("MYSQL_PORT", "3306"),
		User:     getenv("MYSQL_USER", "root"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: getenv("MYSQL_DATABASE", "hotel_booking"),
	}
	if _, err := exec.LookPath("mysql"); err != nil {
		log.Printf("mysql sync disabled: mysql CLI not found: %v", err)
		return nil
	}
	log.Printf("mysql sync enabled on %s:%s/%s", s.Host, s.Port, s.Database)
	return s
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (m *MySQLSync) run(sql string) error {
	args := []string{"-h", m.Host, "-P", m.Port, "-u", m.User}
	if m.Password != "" {
		args = append(args, "-p"+m.Password)
	}
	args = append(args, "-e", sql)
	cmd := exec.Command("mysql", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mysql command failed: %w, output=%s", err, string(out))
	}
	return nil
}

func (m *MySQLSync) InitSchema() {
	sql := fmt.Sprintf(`
CREATE DATABASE IF NOT EXISTS %s;
USE %s;
CREATE TABLE IF NOT EXISTS hotels (
  id INT PRIMARY KEY,
  name VARCHAR(255),
  city VARCHAR(120)
);
CREATE TABLE IF NOT EXISTS rooms (
  id INT PRIMARY KEY,
  hotel_id INT,
  name VARCHAR(255),
  type VARCHAR(50),
  stock INT,
  price_per_night INT,
  capacity INT,
  beds INT
);
CREATE TABLE IF NOT EXISTS bookings (
  id INT PRIMARY KEY,
  room_id INT,
  guests INT,
  nights INT,
  total INT,
  status VARCHAR(40),
  booked_at DATETIME
);`, m.Database, m.Database)
	if err := m.run(sql); err != nil {
		log.Printf("mysql schema init failed: %v", err)
	}
}

func esc(v string) string { return strings.ReplaceAll(v, "'", "''") }

func (m *MySQLSync) SyncHotelsAndRooms(hotels []models.Hotel, rooms []models.Room) {
	for _, h := range hotels {
		sql := fmt.Sprintf("USE %s; REPLACE INTO hotels (id,name,city) VALUES (%d,'%s','%s');", m.Database, h.ID, esc(h.Name), esc(h.City))
		if err := m.run(sql); err != nil {
			log.Printf("mysql hotel sync failed: %v", err)
			break
		}
	}
	for _, r := range rooms {
		sql := fmt.Sprintf("USE %s; REPLACE INTO rooms (id,hotel_id,name,type,stock,price_per_night,capacity,beds) VALUES (%d,%d,'%s','%s',%d,%d,%d,%d);", m.Database, r.ID, r.HotelID, esc(r.Name), esc(string(r.Type)), r.Stock, r.PricePerNight, r.Capacity, r.Beds)
		if err := m.run(sql); err != nil {
			log.Printf("mysql room sync failed: %v", err)
			break
		}
	}
}

func (m *MySQLSync) SyncRoomStock(room models.Room) {
	sql := fmt.Sprintf("USE %s; UPDATE rooms SET stock=%d WHERE id=%d;", m.Database, room.Stock, room.ID)
	if err := m.run(sql); err != nil {
		log.Printf("mysql stock sync failed: %v", err)
	}
}

func (m *MySQLSync) SyncBooking(b models.Booking) {
	sql := fmt.Sprintf("USE %s; REPLACE INTO bookings (id,room_id,guests,nights,total,status,booked_at) VALUES (%d,%d,%d,%d,%d,'%s','%s');", m.Database, b.ID, b.RoomID, b.Guests, b.Nights, b.Total, esc(b.Status), b.BookedAt.Format("2006-01-02 15:04:05"))
	if err := m.run(sql); err != nil {
		log.Printf("mysql booking sync failed: %v", err)
	}
}

func (db *InMemoryDB) SetRoomStock(roomID, stock int) bool {
	if stock < 0 {
		stock = 0
	}
	for i := range db.Rooms {
		if db.Rooms[i].ID != roomID {
			continue
		}
		db.Rooms[i].Stock = stock
		if db.MySQL != nil {
			db.MySQL.SyncRoomStock(db.Rooms[i])
		}
		return true
	}
	return false
}

func (db *InMemoryDB) CreateBooking(roomID, nights, guests int) (models.Room, models.Booking, error) {
	if nights <= 0 {
		nights = 1
	}
	if guests <= 0 {
		guests = 1
	}
	roomIdx := -1
	for i := range db.Rooms {
		if db.Rooms[i].ID == roomID {
			roomIdx = i
			break
		}
	}
	if roomIdx < 0 {
		return models.Room{}, models.Booking{}, fmt.Errorf("kamar tidak ditemukan")
	}
	if db.Rooms[roomIdx].Stock <= 0 {
		return db.Rooms[roomIdx], models.Booking{}, fmt.Errorf("stok kamar habis")
	}
	db.Rooms[roomIdx].Stock--
	now := time.Now()
	booking := models.Booking{
		ID:           len(db.Bookings) + 1,
		UserID:       3,
		RoomID:       roomID,
		Nights:       nights,
		Guests:       guests,
		Total:        db.Rooms[roomIdx].PricePerNight * nights,
		Status:       string(models.BookingConfirmed),
		BookedAt:     now,
		CheckInDate:  now.AddDate(0, 0, 1),
		CheckOutDate: now.AddDate(0, 0, 1+nights),
	}
	db.Bookings = append(db.Bookings, booking)
	if db.MySQL != nil {
		db.MySQL.SyncRoomStock(db.Rooms[roomIdx])
		db.MySQL.SyncBooking(booking)
	}
	return db.Rooms[roomIdx], booking, nil
}

func (db *InMemoryDB) StockByRoomID(roomID int) (int, bool) {
	for _, r := range db.Rooms {
		if r.ID == roomID {
			return r.Stock, true
		}
	}
	return 0, false
}

func ParseIntSafe(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
