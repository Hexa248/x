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
cityCards.forEach(card=>card.addEventListener('click',(e)=>{
  const city=card.dataset.city;
  const image=card.dataset.image;
  const map=card.dataset.map;
  const hero=document.getElementById('city-hero');
  const mapFrame=document.getElementById('main-city-map');
  const homeSearchInput=document.getElementById('homeSearchInput');
  const cityInputHotels=document.getElementById('cityInputHotels');
  if(hero && image){ hero.style.backgroundImage=`linear-gradient(120deg, rgba(14,77,146,.84), rgba(0,169,255,.75)), url('${image}')`; }
  if(mapFrame && map){ mapFrame.src=map; }
  if(homeSearchInput) homeSearchInput.value=city;
  renderCityHotels(city);
  if(cityInputHotels) cityInputHotels.value=city;

  if(!card.classList.contains('active-city-card')){
    e.preventDefault();
    cityCards.forEach(c=>c.classList.remove('active-city-card'));
    card.classList.add('active-city-card');
  }
}));



const homeCityLauncher=document.getElementById('homeCityLauncher');
const homeSearchPopup=document.getElementById('home-search-popup');
const homeSearchInput=document.getElementById('homeSearchInput');
const homeSearchSuggestionList=document.getElementById('homeSearchSuggestionList');
const homeSearchBtn=document.getElementById('homeSearchBtn');
if(homeCityLauncher && homeSearchPopup && homeSearchInput && homeSearchSuggestionList){
  const cityMeta={
    'Bandung':{subtitle:'Jawa Barat, Indonesia',hotels:'4.810 hotel'},
    'Balikpapan':{subtitle:'Kalimantan Timur, Indonesia',hotels:'1.240 hotel'},
    'Banda Aceh':{subtitle:'Aceh, Indonesia',hotels:'560 hotel'},
    'Banjarbaru':{subtitle:'Kalimantan Selatan, Indonesia',hotels:'320 hotel'},
    'Banjarmasin':{subtitle:'Kalimantan Selatan, Indonesia',hotels:'790 hotel'},
    'Batam':{subtitle:'Kepulauan Riau, Indonesia',hotels:'2.050 hotel'}
  };

  const defaultCities=[
    'Bandung','Balikpapan','Banda Aceh','Banjarbaru','Banjarmasin','Batam','Bogor','Cirebon',
    'Denpasar','Jakarta','Jayapura','Kupang','Makassar','Malang','Manado','Medan','Padang',
    'Palembang','Pekanbaru','Pontianak','Semarang','Solo','Surabaya','Yogyakarta'
  ];
  const homeEntries=[];

  cityCards.forEach(card=>{
    const city=(card.dataset.city||'').trim();
    const hotel=(card.dataset.hotel||'').trim();
    if(city) homeEntries.push({label:city,type:'city'});
    if(hotel) homeEntries.push({label:hotel,type:'hotel'});
  });

  defaultCities.forEach(city=>homeEntries.push({label:city,type:'city'}));
  const uniqueEntries=[];
  const seen=new Set();
  homeEntries.forEach(entry=>{
    const key=`${entry.type}:${entry.label.toLowerCase()}`;
    if(!seen.has(key)){
      seen.add(key);
      uniqueEntries.push(entry);
    }
  });

  let selectedEntry={label:'',type:'city'};

  const applyHomeCardFilter=(term)=>{
    const keyword=(term||'').toLowerCase();
    cityCards.forEach(card=>{
      const city=(card.dataset.city||'').toLowerCase();
      const show=!keyword || city.includes(keyword);
      card.classList.toggle('hidden', !show);
    });
  };

  const scoreSuggestion=(item, keyword)=>{
    const text=item.label.toLowerCase();
    if(!keyword) return item.type==='city' ? 0 : item.type==='hotel' ? 1 : 2;
    const starts=text.startsWith(keyword);
    if(!text.includes(keyword)) return 999;
    const typeBoost=item.type==='city' ? 0 : item.type==='hotel' ? 1 : 2;
    return (starts ? 0 : 10) + typeBoost;
  };

  const detailOf=(item)=>{
    if(item.type==='city'){
      const meta=cityMeta[item.label] || {subtitle:'Indonesia',hotels:'Hotel tersedia'};
      return {badge:'Kota', subtitle:meta.subtitle, extra:meta.hotels};
    }
    return {badge:'Hotel', subtitle:'Pilihan properti', extra:'Lihat kamar'};
  };

  const renderHomeSuggestions=(term='')=>{
    const keyword=term.toLowerCase().trim();
    let matched=uniqueEntries
      .filter(item=>item.type==='city')
      .filter(item=>!keyword || item.label.toLowerCase().includes(keyword))
      .sort((a,b)=>{
        const scoreA=scoreSuggestion(a, keyword);
        const scoreB=scoreSuggestion(b, keyword);
        if(scoreA!==scoreB) return scoreA-scoreB;
        return a.label.localeCompare(b.label,'id');
      });

    matched=matched.slice(0,10);

    if(!matched.length){
      homeSearchSuggestionList.innerHTML='<p class="muted">Tidak ada hasil pencarian.</p>';
      return;
    }

    homeSearchSuggestionList.innerHTML=matched.map(item=>{
      const detail=detailOf(item);
      return `
      <button class="suggestion-item suggestion-rich" data-type="${item.type}" data-label="${item.label}" type="button">
        <div>
          <strong>${item.label}</strong>
          <p>${detail.subtitle}</p>
        </div>
        <div class="suggestion-meta">
          <span class="suggestion-badge">${detail.badge}</span>
          <small>${detail.extra}</small>
        </div>
      </button>
    `;
    }).join('');

    homeSearchSuggestionList.querySelectorAll('.suggestion-item').forEach(btn=>{
      btn.addEventListener('click',()=>{
        selectedEntry={label:btn.dataset.label||'', type:btn.dataset.type||'city'};
        homeSearchInput.value=selectedEntry.label;
        applyHomeCardFilter(selectedEntry.label);
        closeHomeSearchPopup();
      });
    });
  };

  const openHomeSearchPopup=()=>{
    homeSearchPopup.classList.remove('hidden');
    renderHomeSuggestions(homeSearchInput.value);
  };

  const closeHomeSearchPopup=()=>homeSearchPopup.classList.add('hidden');

  homeCityLauncher.addEventListener('click', openHomeSearchPopup);
  homeSearchInput.addEventListener('focus', openHomeSearchPopup);

  homeSearchInput.addEventListener('input',(e)=>{
    const value=e.target.value||'';
    selectedEntry={label:value,type:'city'};
    openHomeSearchPopup();
    renderHomeSuggestions(value);
    applyHomeCardFilter(value);
  });

  document.addEventListener('click',(e)=>{
    const target=e.target;
    const insidePanel=homeSearchPopup.contains(target);
    const onLauncher=homeCityLauncher.contains(target);
    if(!insidePanel && !onLauncher) closeHomeSearchPopup();
  });

  homeSearchInput.addEventListener('keydown',(e)=>{
    if(e.key==='Escape') closeHomeSearchPopup();
  });

  homeSearchBtn?.addEventListener('click',()=>{
    const keyword=(homeSearchInput.value || selectedEntry.label || '').trim();
    if(!keyword) return;
    const cityMatch=uniqueEntries.find(item=>item.type==='city' && item.label.toLowerCase()===keyword.toLowerCase());
    if(cityMatch){
      window.location.href=`/hotels?city=${encodeURIComponent(cityMatch.label)}`;
      return;
    }
    window.location.href=`/hotels?q=${encodeURIComponent(keyword)}`;
  });

  renderHomeSuggestions('');
}


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

  const moveIndicator=(animated=true)=>{
    if(!tabsTrack) return;
    const active=tabsTrack.querySelector('.tab-anchor.active');
    if(!active){
      if(indicator) indicator.style.width='0px';
      return;
    }
    if(indicator){
      indicator.style.width=`${active.offsetWidth}px`;
      indicator.style.transform=`translateX(${active.offsetLeft}px)`;
    }

    const targetLeft=active.offsetLeft - (tabsTrack.clientWidth/2) + (active.offsetWidth/2);
    tabsTrack.scrollTo({
      left:Math.max(0,targetLeft),
      behavior:animated ? 'smooth' : 'auto'
    });
  };

  let lastActiveId='';

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
    const changed=lastActiveId!==current.id;
    lastActiveId=current.id;
    moveIndicator(changed);
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
  window.addEventListener('resize', ()=>{ updateDockedState(); moveIndicator(false); updateActiveTab(); });
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
      if(entry.isIntersecting){
        entry.target.classList.add('in-view');
      }else{
        entry.target.classList.remove('in-view');
      }
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
  const smartSearchBtn=document.getElementById('smartSearchBtn');
  const floatSearchBtn=document.getElementById('floatSearchBtn');
  const floatFocusHotel=document.getElementById('floatFocusHotel');
  const hotelPullHandle=document.getElementById('hotelPullHandle');
  const floatingSearchSheet=document.getElementById('floatingSearchSheet');
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
  const mobileFilterBtn=document.getElementById('mobileFilterBtn');
  const filterDrawerBackdrop=document.getElementById('filterDrawerBackdrop');
  const hotelFilterSidebar=document.getElementById('hotelFilterSidebar');
  const priceMin=document.getElementById('priceMin');
  const priceMax=document.getElementById('priceMax');
  const priceMinLabel=document.getElementById('priceMinLabel');
  const priceMaxLabel=document.getElementById('priceMaxLabel');
  const smartDateText=document.getElementById('smartDateText');
  const smartGuestText=document.getElementById('smartGuestText');
  const floatCityText=document.getElementById('floatCityText');
  const floatDateText=document.getElementById('floatDateText');
  const floatGuestText=document.getElementById('floatGuestText');

  const formatIDR=(v)=>`IDR ${Number(v).toLocaleString('id-ID')}`;
  const state={
    city:(cityInputHotels?.value||'').trim(),
    checkIn:'Min, 22 Feb 2026',
    checkOut:'Sen, 23 Feb 2026',
    adults:2,
    children:0,
    rooms:1,
    minPrice:Number(priceMin?.value||100000),
    maxPrice:Number(priceMax?.value||2000000)
  };

  const refreshSearchInputs=()=>{
    const cityText=state.city || 'Pilih Kota';
    const dateText=`${state.checkIn} - ${state.checkOut}`;
    const guestText=`${state.adults} Dewasa, ${state.children} Anak, ${state.rooms} Kamar`;
    if(cityInputHotels) cityInputHotels.value=state.city;
    if(smartDateText) smartDateText.textContent=dateText;
    if(smartGuestText) smartGuestText.textContent=guestText;
    if(floatCityText) floatCityText.textContent=cityText || 'Cari nama hotel';
    if(floatDateText) floatDateText.textContent=dateText;
    if(floatGuestText) floatGuestText.textContent=guestText;
  };

  const openPopup=(kind)=>{
    if(!popup || !popupBody || !popupTitle) return;
    popup.dataset.kind=kind;
    popupBody.innerHTML='';

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

  const tagFilters=(key)=>[...document.querySelectorAll(`input[data-filter="${key}"]:checked`)].map(x=>x.value);

  const accommodationTypes=['hotel','vila','apartemen','guesthouse','ryokan'];
  const promoTags=['promo-anda','ramadan','extra-benefit','domestik'];
  const facilityTags=['restoran','sarapan','early-checkin','antar-jemput-bandara','late-checkout'];
  cards.forEach((card,idx)=>{
    card.dataset.type=accommodationTypes[idx%accommodationTypes.length];
    const promoA=promoTags[idx%promoTags.length];
    const promoB=promoTags[(idx+2)%promoTags.length];
    const facilityA=facilityTags[idx%facilityTags.length];
    const facilityB=facilityTags[(idx+1)%facilityTags.length];
    card.dataset.promos=`${promoA},${promoB}`;
    card.dataset.facilities=`${facilityA},${facilityB}`;
  });

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
    const types=tagFilters('type');
    const promos=tagFilters('promo');
    const facilities=tagFilters('facility');
    const keyword=(state.city||'').toLowerCase();

    let visible=0;
    cards.forEach(card=>{
      const city=(card.dataset.city||'').toLowerCase();
      const name=(card.dataset.name||'').toLowerCase();
      const rating=Number(card.dataset.rating||0);
      const star=Number(card.dataset.star||0);
      const price=Number(card.dataset.price||0);
      const type=(card.dataset.type||'').toLowerCase();
      const cardPromos=(card.dataset.promos||'').split(',').filter(Boolean);
      const cardFacilities=(card.dataset.facilities||'').split(',').filter(Boolean);

      const keywordMatch=!keyword || name.includes(keyword) || city.includes(keyword);
      const starMatch=!stars.length || stars.includes(star);
      const ratingMatch=!ratings.length || ratings.some(r=>rating>=r);
      const priceMatch=price>=state.minPrice && price<=state.maxPrice;
      const typeMatch=!types.length || types.includes(type);
      const promoMatch=!promos.length || promos.some(tag=>cardPromos.includes(tag));
      const facilityMatch=!facilities.length || facilities.some(tag=>cardFacilities.includes(tag));

      const show=keywordMatch && starMatch && ratingMatch && priceMatch && typeMatch && promoMatch && facilityMatch;
      card.classList.toggle('hidden', !show);
      if(show) visible+=1;
    });

    if(summary){
      summary.textContent=visible>0 ? `Menampilkan ${visible} hotel sesuai pencarian & filter aktif.` : 'Tidak ada hasil, coba ubah filter atau kata kunci.';
    }
    if(emptyState) emptyState.classList.toggle('hidden', visible>0);
  };

  popupCancel?.addEventListener('click', closePopup);
  popup?.addEventListener('click',(e)=>{ if(e.target===popup) closePopup(); });

  document.querySelectorAll('[data-open]').forEach(el=>{
    el.addEventListener('click',()=>openPopup(el.dataset.open));
  });

  cityInputHotels?.addEventListener('input',(e)=>{
    state.city=(e.target.value||'').trim();
    if(floatCityText) floatCityText.textContent=state.city || 'Cari nama hotel';
    applyFilters();
  });

  floatFocusHotel?.addEventListener('click',()=>{
    cityInputHotels?.focus();
    cityInputHotels?.scrollIntoView({behavior:'smooth', block:'center'});
  });

  popupApply?.addEventListener('click',()=>{
    const kind=popup?.dataset.kind;
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

  const openFilterDrawer=()=>{ document.body.classList.add('filter-drawer-open'); filterDrawerBackdrop?.classList.remove('hidden'); };
  const closeFilterDrawer=()=>{ document.body.classList.remove('filter-drawer-open'); filterDrawerBackdrop?.classList.add('hidden'); };
  mobileFilterBtn?.addEventListener('click', openFilterDrawer);
  filterDrawerBackdrop?.addEventListener('click', closeFilterDrawer);

  let touchStartX=0;
  let touchCurrentX=0;
  let fromEdge=false;

  window.addEventListener('touchstart',(e)=>{
    touchStartX=e.touches[0].clientX;
    touchCurrentX=touchStartX;
    fromEdge=touchStartX < 24;
  },{passive:true});

  window.addEventListener('touchmove',(e)=>{ touchCurrentX=e.touches[0].clientX; },{passive:true});

  window.addEventListener('touchend',()=>{
    const delta=touchCurrentX-touchStartX;
    if(window.innerWidth>900) return;
    if(fromEdge && delta>70) openFilterDrawer();
    if(document.body.classList.contains('filter-drawer-open') && delta<-60) closeFilterDrawer();
  },{passive:true});

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
    lastScrollY=window.scrollY;
    const activePane=document.querySelector('.split-scroll-right') || document.querySelector('.split-scroll-left');
    lastPaneY=activePane ? activePane.scrollTop : 0;
    popupShift=0;
    const content=filterActionPopup.querySelector('.filter-action-popup-content');
    if(content) content.style.transform='translateY(0)';
  };

  let lastScrollY=window.scrollY;
  let lastPaneY=0;
  let popupShift=0;
  const maxPopupShift=120;

  const updateFilterPopupShift=(delta)=>{
    if(!filterActionPopup || filterActionPopup.classList.contains('hidden')) return;
    const content=filterActionPopup.querySelector('.filter-action-popup-content');
    if(!content) return;
    popupShift=Math.max(0, Math.min(maxPopupShift, popupShift+delta));
    content.style.transform=`translateY(${popupShift}px)`;
  };

  const handleWindowScroll=()=>{
    const currentY=window.scrollY;
    const delta=currentY-lastScrollY;
    updateFilterPopupShift(delta);
    lastScrollY=currentY;
  };

  const handlePaneScroll=(e)=>{
    const pane=e.currentTarget;
    const current=pane.scrollTop;
    const delta=current-lastPaneY;
    updateFilterPopupShift(delta);
    lastPaneY=current;
  };

  const closeFilterActionPopup=()=>{
    if(!filterActionPopup) return;
    filterActionPopup.classList.add('hidden');
    popupShift=0;
    const content=filterActionPopup.querySelector('.filter-action-popup-content');
    if(content) content.style.transform='translateY(0)';
  };

  window.addEventListener('scroll', handleWindowScroll, {passive:true});
  document.querySelectorAll('.split-scroll-pane').forEach(p=>p.addEventListener('scroll', handlePaneScroll, {passive:true}));

  [priceMin, priceMax].forEach(el=>el?.addEventListener('input',applyFilters));
  document.querySelectorAll('input[data-filter]').forEach(el=>el.addEventListener('change',applyFilters));
  smartSearchBtn?.addEventListener('click',applyFilters);
  floatSearchBtn?.addEventListener('click',applyFilters);
  applyFiltersBtn?.addEventListener('click',()=>openFilterActionPopup('apply'));
  resetFiltersBtn?.addEventListener('click',()=>openFilterActionPopup('reset'));

  filterActionCancel?.addEventListener('click', closeFilterActionPopup);
  filterActionPopup?.addEventListener('click',(e)=>{ if(e.target===filterActionPopup) closeFilterActionPopup(); });
  filterActionConfirm?.addEventListener('click',()=>{
    const action=filterActionPopup?.dataset.action;
    if(action==='reset') doResetFilters();
    if(action==='apply') applyFilters();
    closeFilterActionPopup();
    closeFilterDrawer();
  });

  let lastWindowY=window.scrollY;
  const updateTopbarAndPull=()=>{
    const topbar=document.querySelector('.topbar');
    if(!topbar) return;
    const y=window.scrollY;
    const nearTop=y < 24;
    const scrollingDown=y > lastWindowY;
    if(nearTop || !scrollingDown){
      document.body.classList.remove('topbar-hidden');
      hotelPullHandle?.classList.add('hidden');
      floatingSearchSheet?.classList.add('hidden');
    }else{
      document.body.classList.add('topbar-hidden');
      hotelPullHandle?.classList.remove('hidden');
    }
    lastWindowY=y;
  };

  hotelPullHandle?.addEventListener('click',()=>{
    floatingSearchSheet?.classList.toggle('hidden');
  });

  window.addEventListener('scroll', updateTopbarAndPull, {passive:true});
  window.addEventListener('resize',()=>{ if(window.innerWidth>900) closeFilterDrawer(); });
  updateTopbarAndPull();

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

const roomSearchInput=document.getElementById('roomSearchInput');
const roomSearchBtn=document.getElementById('roomSearchBtn');
if(roomSearchInput){
  const roomRows=[...document.querySelectorAll('.room-row')];
  const applyRoomSearch=()=>{
    const keyword=(roomSearchInput.value||'').toLowerCase().trim();
    roomRows.forEach(row=>{
      const name=(row.dataset.name||'').toLowerCase();
      const capacity=row.dataset.capacity||'';
      const beds=row.dataset.beds||'';
      const text=`${name} ${capacity} ${beds}`;
      const show=!keyword || text.includes(keyword);
      row.classList.toggle('hidden', !show);
    });
  };
  roomSearchInput.addEventListener('input', applyRoomSearch);
  roomSearchBtn?.addEventListener('click', applyRoomSearch);
}
