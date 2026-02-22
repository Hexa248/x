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
  if(cityInputHotels) cityInputHotels.value=`${city}, Indonesia`;
}));

const sections=document.querySelectorAll('.info-section');
const tabs=document.querySelectorAll('.tab-anchor');
if(sections.length && tabs.length){
  const updateActiveTab=()=>{
    const sticky=document.querySelector('.sticky-tabs');
    const headerOffset=(sticky ? sticky.offsetHeight : 0) + 70;
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
  };

  window.addEventListener('scroll', updateActiveTab, {passive:true});
  window.addEventListener('resize', updateActiveTab);
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
