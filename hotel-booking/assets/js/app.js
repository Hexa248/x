function openTab(id){
  document.querySelectorAll('.room-tab').forEach(x=>x.classList.add('hidden'));
  const target=document.getElementById(id);
  if(target)target.classList.remove('hidden');
}

function closePopup(){
  const p=document.getElementById('hotel-popup');
  if(p)p.classList.add('hidden');
}

const cards=document.querySelectorAll('.hotel-card');
cards.forEach(c=>c.addEventListener('click',(e)=>{
  if(e.target.closest('a')) return;
  const p=document.getElementById('hotel-popup');
  if(!p) return;
  const name=c.dataset.hotel||'hotel pilihanmu';
  const text=document.getElementById('popup-text');
  if(text) text.innerText='Kamu memilih '+name;
  p.classList.remove('hidden');
}));

const cityCards=document.querySelectorAll('.city-card');
cityCards.forEach(card=>card.addEventListener('click',()=>{
  const city=card.dataset.city;
  const image=card.dataset.image;
  const map=card.dataset.map;
  const hero=document.getElementById('city-hero');
  const mapFrame=document.getElementById('main-city-map');
  const cityInput=document.getElementById('cityInput');
  const cityInputHotels=document.getElementById('cityInputHotels');
  if(hero && image){ hero.style.backgroundImage=`linear-gradient(120deg, rgba(14,77,146,.84), rgba(0,169,255,.75)), url('${image}')`; }
  if(mapFrame && map){ mapFrame.src=map; }
  if(cityInput) cityInput.value=`${city}, Indonesia`;
  renderCityHotels(city);
  if(cityInputHotels) cityInputHotels.value=`${city}, Indonesia`;
}));

const sections=document.querySelectorAll('.info-section');
const tabs=document.querySelectorAll('.tab-anchor');
if(sections.length && tabs.length){
  const sticky=document.querySelector('.sticky-tabs');
  const tabsTrack=document.querySelector('.sticky-tabs-links');
  const indicator=document.querySelector('.tab-indicator');

  const updateDockedState=()=>{
    if(!sticky) return;
    const trigger=(sticky.parentElement ? sticky.parentElement.offsetTop : 0) + 40;
    sticky.classList.toggle('is-docked', window.scrollY > trigger);
  };

  const moveIndicator=()=>{
    if(!indicator || !tabsTrack) return;
    const active=tabsTrack.querySelector('.tab-anchor.active');
    if(!active){
      indicator.style.width='0px';
      return;
    }
    indicator.style.width=`${active.offsetWidth}px`;
    indicator.style.transform=`translateX(${active.offsetLeft}px)`;
  };

  const updateActiveTab=()=>{
    const headerOffset=(sticky ? sticky.offsetHeight : 0) + 14;
    const y=window.scrollY + headerOffset;

    let current=sections[0];
    sections.forEach(sec=>{
      if(y >= sec.offsetTop){
        current=sec;
      }
    });

    tabs.forEach(t=>t.classList.remove('active'));
    const active=document.querySelector(`.tab-anchor[href="#${current.id}"]`);
    if(active) active.classList.add('active');
    moveIndicator();
  };

  const smoothTargets=document.querySelectorAll('.tab-anchor, .sticky-tabs-actions a[href^="#"]');
  smoothTargets.forEach(link=>{
    link.addEventListener('click',(e)=>{
      const href=link.getAttribute('href') || '';
      if(!href.startsWith('#') || href === '#') return;
      const target=document.querySelector(href);
      if(!target) return;
      e.preventDefault();
      const offset=(sticky ? sticky.offsetHeight : 0) + 10;
      const targetY=target.getBoundingClientRect().top + window.scrollY - offset;
      window.scrollTo({top:targetY, behavior:'smooth'});
    });
  });

  window.addEventListener('scroll', ()=>{ updateDockedState(); updateActiveTab(); }, {passive:true});
  window.addEventListener('resize', ()=>{ updateDockedState(); moveIndicator(); updateActiveTab(); });
  updateDockedState();
  updateActiveTab();
}


const recoCards=document.querySelectorAll('.reco-card');
const recoLinks=document.querySelectorAll('.reco-link');
if(recoCards.length && recoLinks.length){
  const recoObserver=new IntersectionObserver((entries)=>{
    entries.forEach(entry=>{
      if(entry.isIntersecting){
        recoLinks.forEach(l=>l.classList.remove('active'));
        const id=entry.target.getAttribute('id');
        const active=document.querySelector(`.reco-link[href="#${id}"]`);
        if(active) active.classList.add('active');
      }
    });
  },{threshold:0.45});
  recoCards.forEach(c=>recoObserver.observe(c));
}


document.body.classList.add('js-animate');
const animatedEls=document.querySelectorAll('.scroll-animate');
if(animatedEls.length){
  const animObserver=new IntersectionObserver((entries)=>{
    entries.forEach(entry=>{
      if(entry.isIntersecting) entry.target.classList.add('in-view');
    });
  },{threshold:0.18});
  animatedEls.forEach(el=>animObserver.observe(el));
}


const cityHotelCatalog={
  "Jakarta":[
    {name:"Jakarta Grand Central",rating:"9.0",price:"Rp600000",image:"https://images.unsplash.com/photo-1566073771259-6a8506099945?auto=format&fit=crop&w=900&q=80"},
    {name:"Sudirman Urban Suites",rating:"8.8",price:"Rp520000",image:"https://images.unsplash.com/photo-1455587734955-081b22074882?auto=format&fit=crop&w=900&q=80"},
    {name:"Menteng Royal Inn",rating:"8.7",price:"Rp470000",image:"https://images.unsplash.com/photo-1591088398332-8a7791972843?auto=format&fit=crop&w=900&q=80"},
    {name:"Kemang Park Hotel",rating:"8.6",price:"Rp450000",image:"https://images.unsplash.com/photo-1445019980597-93fa8acb246c?auto=format&fit=crop&w=900&q=80"},
    {name:"Ancol Bay Resort",rating:"8.9",price:"Rp580000",image:"https://images.unsplash.com/photo-1496417263034-38ec4f0b665a?auto=format&fit=crop&w=900&q=80"}
  ],
  "Bandung":[
    {name:"Bandung Sky Inn",rating:"8.9",price:"Rp420000",image:"https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?auto=format&fit=crop&w=900&q=80"},
    {name:"Cihampelas Urban Stay",rating:"8.7",price:"Rp390000",image:"https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=900&q=80"},
    {name:"Dago Hills Hotel",rating:"8.8",price:"Rp460000",image:"https://images.unsplash.com/photo-1564501049412-61c2a3083791?auto=format&fit=crop&w=900&q=80"},
    {name:"Braga Heritage Inn",rating:"8.6",price:"Rp410000",image:"https://images.unsplash.com/photo-1578683010236-d716f9a3f461?auto=format&fit=crop&w=900&q=80"},
    {name:"Lembang Valley Resort",rating:"9.1",price:"Rp590000",image:"https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd?auto=format&fit=crop&w=900&q=80"}
  ],
  "Surabaya":[
    {name:"Tunjungan City Hotel",rating:"8.8",price:"Rp480000",image:"https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?auto=format&fit=crop&w=900&q=80"},
    {name:"Pakuwon Suites",rating:"8.7",price:"Rp450000",image:"https://images.unsplash.com/photo-1562438668-bcf0ca6578f0?auto=format&fit=crop&w=900&q=80"},
    {name:"Manyar Prime Stay",rating:"8.6",price:"Rp410000",image:"https://images.unsplash.com/photo-1568495248636-6432b97bd949?auto=format&fit=crop&w=900&q=80"},
    {name:"Kenjeran Bay Hotel",rating:"8.5",price:"Rp390000",image:"https://images.unsplash.com/photo-1519821172141-b5d8dbb7db2b?auto=format&fit=crop&w=900&q=80"},
    {name:"Surabaya Grand Palace",rating:"9.0",price:"Rp620000",image:"https://images.unsplash.com/photo-1578898887932-dce23a595ad4?auto=format&fit=crop&w=900&q=80"}
  ],
  "Yogyakarta":[
    {name:"Yogyakarta Heritage Stay",rating:"8.7",price:"Rp350000",image:"https://images.unsplash.com/photo-1521783593447-5702b9bfd267?auto=format&fit=crop&w=900&q=80"},
    {name:"Malioboro City Inn",rating:"8.8",price:"Rp420000",image:"https://images.unsplash.com/photo-1613977257592-487ecd136cc3?auto=format&fit=crop&w=900&q=80"},
    {name:"Tugu Art Hotel",rating:"8.9",price:"Rp480000",image:"https://images.unsplash.com/photo-1468824357306-a439d58ccb1c?auto=format&fit=crop&w=900&q=80"},
    {name:"Kaliurang Breeze",rating:"8.5",price:"Rp330000",image:"https://images.unsplash.com/photo-1568084680786-a84f91d1153c?auto=format&fit=crop&w=900&q=80"},
    {name:"Jogja Royal Resort",rating:"9.1",price:"Rp610000",image:"https://images.unsplash.com/photo-1520250497591-112f2f40a3f4?auto=format&fit=crop&w=900&q=80"}
  ]
};

function renderCityHotels(city){
  const grid=document.getElementById('cityHotelGrid');
  const title=document.getElementById('cityHotelTitle');
  const sub=document.getElementById('cityHotelSub');
  if(!grid||!title||!sub) return;

  const source=cityHotelCatalog[city] || [
    {name:`${city} Vista Resort`,rating:'8.8',price:'Rp450000',image:'https://images.unsplash.com/photo-1445019980597-93fa8acb246c?auto=format&fit=crop&w=900&q=80'},
    {name:`${city} City Inn`,rating:'8.6',price:'Rp390000',image:'https://images.unsplash.com/photo-1540541338287-41700207dee6?auto=format&fit=crop&w=900&q=80'},
    {name:`${city} Grand Central`,rating:'8.9',price:'Rp520000',image:'https://images.unsplash.com/photo-1455587734955-081b22074882?auto=format&fit=crop&w=900&q=80'},
    {name:`${city} Heritage Stay`,rating:'8.7',price:'Rp430000',image:'https://images.unsplash.com/photo-1496417263034-38ec4f0b665a?auto=format&fit=crop&w=900&q=80'},
    {name:`${city} Sky Suites`,rating:'9.0',price:'Rp580000',image:'https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=900&q=80'}
  ];

  title.textContent=`5 Hotel Pilihan di ${city}`;
  sub.textContent=`Rekomendasi hotel terbaik di ${city}, desain sudah disesuaikan dengan referensi.`;
  grid.innerHTML=source.slice(0,5).map(h=>`
    <article class="city-hotel-card">
      <img src="${h.image}" alt="${h.name}">
      <div class="card-body">
        <h4>${h.name}</h4>
        <p>${city} • Rating ${h.rating}</p>
        <p>Mulai dari <strong>${h.price}/malam</strong></p>
        <a class="btn" href="/hotels">Lihat Detail</a>
      </div>
    </article>
  `).join('');
}




const hotelListing=document.getElementById('hotelListing');
if(hotelListing){
  const cityInputHotels=document.getElementById('cityInputHotels');
  const dateInputHotels=document.getElementById('dateInputHotels');
  const guestInputHotels=document.getElementById('guestInputHotels');
  const searchHotelsBtn=document.getElementById('searchHotelsBtn');
  const cards=[...hotelListing.querySelectorAll('.result-card')];
  const summary=document.getElementById('hotelResultSummary');
  const emptyState=document.getElementById('hotelEmptyState');
  const applyFiltersBtn=document.getElementById('applyFiltersBtn');
  const resetFiltersBtn=document.getElementById('resetFiltersBtn');
  const popup=document.getElementById('hotel-search-popup');
  const popupTitle=document.getElementById('hotelPopupTitle');
  const popupBody=document.getElementById('hotelPopupBody');
  const popupApply=document.getElementById('hotelPopupApply');
  const popupCancel=document.getElementById('hotelPopupCancel');
  const filterActionPopup=document.getElementById('filter-action-popup');
  const filterActionTitle=document.getElementById('filterActionTitle');
  const filterActionText=document.getElementById('filterActionText');
  const filterActionCancel=document.getElementById('filterActionCancel');
  const filterActionConfirm=document.getElementById('filterActionConfirm');
  const priceMin=document.getElementById('priceMin');
  const priceMax=document.getElementById('priceMax');
  const priceMinLabel=document.getElementById('priceMinLabel');
  const priceMaxLabel=document.getElementById('priceMaxLabel');

  const formatIDR=(v)=>`IDR ${Number(v).toLocaleString('id-ID')}`;
  const state={
    city:(cityInputHotels?.value||'').replace(', Indonesia','').trim(),
    checkIn:'Min, 22 Feb 2026',
    checkOut:'Sen, 23 Feb 2026',
    adults:2,
    children:0,
    rooms:1,
    minPrice:Number(priceMin?.value||100000),
    maxPrice:Number(priceMax?.value||2000000)
  };

  const refreshSearchInputs=()=>{
    if(cityInputHotels) cityInputHotels.value=state.city ? `${state.city}, Indonesia` : '';
    if(dateInputHotels) dateInputHotels.value=`${state.checkIn} - ${state.checkOut}`;
    if(guestInputHotels) guestInputHotels.value=`${state.adults} Dewasa, ${state.children} Anak, ${state.rooms} Kamar`;
  };

  const openPopup=(kind)=>{
    if(!popup || !popupBody || !popupTitle) return;
    popup.dataset.kind=kind;
    popupBody.innerHTML='';

    if(kind==='city'){
      popupTitle.textContent='Pilih Kota atau Nama Hotel';
      popupBody.innerHTML=`<input id="popupCityValue" placeholder="Contoh: Jakarta atau Nusantara" value="${state.city}" />`;
    }

    if(kind==='date'){
      popupTitle.textContent='Pilih Tanggal Menginap';
      popupBody.innerHTML=`
        <label>Check-in</label>
        <input id="popupCheckIn" type="date" value="2026-02-22" />
        <label>Check-out</label>
        <input id="popupCheckOut" type="date" value="2026-02-23" />
      `;
    }

    if(kind==='guest'){
      popupTitle.textContent='Atur Tamu & Kamar';
      popupBody.innerHTML=`
        <label>Dewasa</label><input id="popupAdults" type="number" min="1" value="${state.adults}" />
        <label>Anak</label><input id="popupChildren" type="number" min="0" value="${state.children}" />
        <label>Kamar</label><input id="popupRooms" type="number" min="1" value="${state.rooms}" />
      `;
    }

    popup.classList.remove('hidden');
  };

  const closePopup=()=>{ if(popup) popup.classList.add('hidden'); };

  const selectedStars=()=>[...document.querySelectorAll('input[data-filter="star"]:checked')].map(x=>Number(x.value));
  const selectedRatings=()=>[...document.querySelectorAll('input[data-filter="rating"]:checked')].map(x=>Number(x.value));

  const applyFilters=()=>{
    if(priceMin && priceMax){
      state.minPrice=Number(priceMin.value);
      state.maxPrice=Number(priceMax.value);
      if(state.minPrice > state.maxPrice){
        const tmp=state.minPrice; state.minPrice=state.maxPrice; state.maxPrice=tmp;
      }
      priceMin.value=String(state.minPrice);
      priceMax.value=String(state.maxPrice);
      if(priceMinLabel) priceMinLabel.textContent=formatIDR(state.minPrice);
      if(priceMaxLabel) priceMaxLabel.textContent=formatIDR(state.maxPrice);
    }

    const stars=selectedStars();
    const ratings=selectedRatings();
    const keyword=(state.city||'').toLowerCase();

    let visible=0;
    cards.forEach(card=>{
      const city=(card.dataset.city||'').toLowerCase();
      const name=(card.dataset.name||'').toLowerCase();
      const rating=Number(card.dataset.rating||0);
      const star=Number(card.dataset.star||0);
      const price=Number(card.dataset.price||0);

      const cityMatch=!keyword || city.includes(keyword) || name.includes(keyword);
      const starMatch=!stars.length || stars.includes(star);
      const ratingMatch=!ratings.length || ratings.some(r=>rating>=r);
      const priceMatch=price>=state.minPrice && price<=state.maxPrice;

      const show=cityMatch && starMatch && ratingMatch && priceMatch;
      card.classList.toggle('hidden', !show);
      if(show) visible+=1;
    });

    if(summary){
      summary.textContent=visible>0 ? `Menampilkan ${visible} hotel sesuai pencarian & filter aktif.` : 'Tidak ada hasil, coba ubah filter atau kata kunci.';
    }
    if(emptyState) emptyState.classList.toggle('hidden', visible>0);
  };

  cityInputHotels?.addEventListener('click',()=>openPopup('city'));
  dateInputHotels?.addEventListener('click',()=>openPopup('date'));
  guestInputHotels?.addEventListener('click',()=>openPopup('guest'));
  popupCancel?.addEventListener('click', closePopup);
  popup?.addEventListener('click',(e)=>{ if(e.target===popup) closePopup(); });

  popupApply?.addEventListener('click',()=>{
    const kind=popup?.dataset.kind;
    if(kind==='city'){
      const value=document.getElementById('popupCityValue');
      state.city=(value?.value||'').trim();
    }
    if(kind==='date'){
      const ci=document.getElementById('popupCheckIn');
      const co=document.getElementById('popupCheckOut');
      if(ci?.value) state.checkIn=new Date(ci.value).toLocaleDateString('id-ID',{weekday:'short', day:'2-digit', month:'short', year:'numeric'});
      if(co?.value) state.checkOut=new Date(co.value).toLocaleDateString('id-ID',{weekday:'short', day:'2-digit', month:'short', year:'numeric'});
    }
    if(kind==='guest'){
      const adults=document.getElementById('popupAdults');
      const children=document.getElementById('popupChildren');
      const rooms=document.getElementById('popupRooms');
      state.adults=Math.max(1, Number(adults?.value||2));
      state.children=Math.max(0, Number(children?.value||0));
      state.rooms=Math.max(1, Number(rooms?.value||1));
    }

    refreshSearchInputs();
    closePopup();
    applyFilters();
  });

  document.querySelectorAll('.quick-date').forEach(btn=>{
    btn.addEventListener('click',()=>{
      const start=new Date();
      start.setDate(start.getDate()+Number(btn.dataset.days||0));
      const end=new Date(start);
      end.setDate(end.getDate()+1);
      state.checkIn=start.toLocaleDateString('id-ID',{weekday:'short', day:'2-digit', month:'short', year:'numeric'});
      state.checkOut=end.toLocaleDateString('id-ID',{weekday:'short', day:'2-digit', month:'short', year:'numeric'});
      refreshSearchInputs();
    });
  });

  const doResetFilters=()=>{
    state.city='';
    state.minPrice=100000;
    state.maxPrice=2000000;
    if(priceMin) priceMin.value='100000';
    if(priceMax) priceMax.value='2000000';
    document.querySelectorAll('input[data-filter]').forEach(chk=>{ chk.checked=false; });
    document.querySelectorAll('input[data-filter="star"]').forEach(chk=>{ if(['3','4','5'].includes(chk.value)) chk.checked=true; });
    refreshSearchInputs();
    applyFilters();
  };

  const openFilterActionPopup=(kind)=>{
    if(!filterActionPopup || !filterActionTitle || !filterActionText) return;
    filterActionPopup.dataset.action=kind;
    if(kind==='apply'){
      filterActionTitle.textContent='Terapkan Filter';
      filterActionText.textContent='Gunakan filter saat ini untuk memperbarui daftar hotel?';
    }else{
      filterActionTitle.textContent='Reset Filter';
      filterActionText.textContent='Reset semua filter ke kondisi default?';
    }
    filterActionPopup.classList.remove('hidden');
  };

  const closeFilterActionPopup=()=>{ if(filterActionPopup) filterActionPopup.classList.add('hidden'); };

  [priceMin, priceMax].forEach(el=>el?.addEventListener('input',applyFilters));
  document.querySelectorAll('input[data-filter]').forEach(el=>el.addEventListener('change',applyFilters));
  searchHotelsBtn?.addEventListener('click',applyFilters);
  applyFiltersBtn?.addEventListener('click',()=>openFilterActionPopup('apply'));
  resetFiltersBtn?.addEventListener('click',()=>openFilterActionPopup('reset'));

  filterActionCancel?.addEventListener('click', closeFilterActionPopup);
  filterActionPopup?.addEventListener('click',(e)=>{ if(e.target===filterActionPopup) closeFilterActionPopup(); });
  filterActionConfirm?.addEventListener('click',()=>{
    const action=filterActionPopup?.dataset.action;
    if(action==='reset') doResetFilters();
    if(action==='apply') applyFilters();
    closeFilterActionPopup();
  });

  refreshSearchInputs();
  applyFilters();
}

const paymentGrid=document.getElementById('paymentMethodGrid');
if(paymentGrid){
  const cards=[...paymentGrid.querySelectorAll('.payment-method-card')];
  const selectedPaymentText=document.getElementById('selectedPaymentText');
  const payNowBtn=document.getElementById('payNowBtn');

  const setActive=(card)=>{
    cards.forEach(c=>c.classList.remove('active'));
    card.classList.add('active');
    const method=card.dataset.method || 'Virtual Account';
    if(selectedPaymentText) selectedPaymentText.innerHTML=`Metode dipilih: <strong>${method}</strong>`;
    if(payNowBtn) payNowBtn.textContent=`Bayar dengan ${method}`;
  };

  cards.forEach(card=>{
    card.addEventListener('click',()=>{
      const radio=card.querySelector('input[type="radio"]');
      if(radio) radio.checked=true;
      setActive(card);
    });
  });

  if(payNowBtn){
    payNowBtn.addEventListener('click',()=>{
      const active=paymentGrid.querySelector('.payment-method-card.active');
      const method=active ? active.dataset.method : 'Virtual Account';
      alert(`Mock pembayaran berhasil diproses lewat ${method}.`);
    });
  }
}
