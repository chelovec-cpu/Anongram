import React, { useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';

const avatar = '/anongram-avatar.png';

type Chat = {id:string; title:string; username:string; preview:string; time:string; unread:number; type:string; avatar:string; online?:boolean};
const chats:Chat[] = [
 {id:'1',title:'Анонимный чат',username:'@anonymous',preview:'Добро пожаловать в AnonGram',time:'12:40',unread:3,type:'Группы',avatar:'A',online:true},
 {id:'2',title:'Hailendsky',username:'@hailendsky',preview:'Увидимся позже',time:'11:22',unread:0,type:'Все чаты',avatar:'H',online:true},
 {id:'3',title:'CRMP MOBILE',username:'@crmpmobile',preview:'Новый релиз уже близко',time:'09:17',unread:7,type:'Каналы',avatar:'C'},
 {id:'4',title:'Избранное',username:'saved',preview:'Сохранённые сообщения',time:'вчера',unread:0,type:'Избранное',avatar:'★'},
 {id:'5',title:'AnonGram News',username:'@anongramnews',preview:'Обновление интерфейса',time:'вчера',unread:12,type:'Каналы',avatar:'N'},
];

function Logo({small=false}:{small?:boolean}) { return <div className={small?'brandMini':'brand'}><img src={avatar}/><div><b>ANONGRAM</b>{!small&&<span>Познай анонимность_</span>}</div></div> }
function Avatar({c,size='normal'}:{c:{avatar:string;online?:boolean};size?:string}) { return <div className={'avatar '+size}>{c.avatar==='A'?<img src={avatar}/>:c.avatar}<i className={c.online?'online':''}/></div> }

function Auth({onLogin}:{onLogin:()=>void}) {
 const [register,setRegister]=useState(false); const [show,setShow]=useState(false);
 return <div className="authPage"><div className="authHole"/><div className="authStars"/><main className="authWrap">
  <Logo/>
  <section className="authCard">
   {register ? <>
    <div className="authTitle">Создать аккаунт</div><div className="authSub">Новый профиль AnonGram</div>
    <label>Имя пользователя<input placeholder="@username"/></label><label>Номер / email<input placeholder="Номер или email"/></label>
    <label>Пароль<div className="password"><input type={show?'text':'password'} placeholder="Пароль"/><button onClick={()=>setShow(!show)}>{show?'Скрыть':'Показать'}</button></div></label>
    <button className="primary" onClick={onLogin}>Зарегистрироваться</button><button className="linkBtn" onClick={()=>setRegister(false)}>Войти в аккаунт</button>
   </> : <>
    <div className="authTitle">Вход в AnonGram</div><div className="authSub">Твоя анонимность начинается здесь</div>
    <label>Логин / Номер<input placeholder="Логин или номер"/></label><label>Пароль<div className="password"><input type={show?'text':'password'} placeholder="Пароль"/><button onClick={()=>setShow(!show)}>{show?'Скрыть':'Показать'}</button></div></label>
    <button className="primary" onClick={onLogin}>Войти</button><div className="or"><span/>или<span/></div><button className="secondary" onClick={()=>setRegister(true)}>Создать аккаунт</button><button className="linkBtn">Забыли пароль?</button>
   </>}
   <small>Продолжая, вы принимаете условия использования AnonGram.</small>
  </section>
 </main></div>
}

function Drawer({close,onProfile,logout}:{close:()=>void;onProfile:()=>void;logout:()=>void}) {
 const noop=()=>{}; const rows: Array<[string,string,()=>void]>=[['👤','Мой профиль',onProfile],['⇄','Аккаунты',noop],['▱','Архив',noop],['★','Звёзды',noop],['◈','TON',noop],['♧','Контакты',noop],['◉','Звонки',noop],['▣','Сохранённые',noop],['◎','Истории',noop],['◐','Темы',noop],['⌕','Глобальный поиск',noop],['♟','Боты',noop]];
 return <><div className="overlay" onClick={close}/><aside className="drawer"><div className="drawerTop"><Logo small/><button onClick={close}>×</button></div><div className="drawerProfile" onClick={onProfile}><img src={avatar}/><b>Hailendsky</b><span>@hailendsky · Founder</span><em>● онлайн</em></div>{rows.map(([i,t,fn])=><button className="drawerRow" key={t} onClick={fn}>{i}<span>{t}</span><small>›</small></button>)}<div className="drawerSep"/><button className="drawerRow"><span className="createDot">+</span><span>Создать чат</span></button><button className="drawerRow"><span className="createDot">+</span><span>Создать группу</span></button><button className="drawerRow"><span className="createDot">+</span><span>Создать канал</span></button><button className="drawerRow"><span>⚙</span><span>Настройки</span></button><button className="drawerRow danger" onClick={logout}><span>↪</span><span>Выйти</span></button></aside></>
}

function Profile({back}:{back:()=>void}) {
 return <main className="profilePage"><div className="profileTop"><button className="roundBtn" onClick={back}>‹</button><b>Профиль</b><button className="roundBtn">⋮</button></div><section className="profileBanner"><div className="bannerHole"/><div className="bannerGlow"/><div className="profileAvatar"><img src={avatar}/><i/></div></section><section className="profileBody"><div className="profileName"><h1>Hailendsky <span>✓</span></h1><p>@hailendsky</p><div className="badges"><b>♛ FOUNDER</b><b>◆ ADMIN</b><b>● Хороший рейтинг</b></div></div><div className="profileActions"><button>Изменить</button><button>Фото</button><button>Настройки</button><button>Ещё</button></div><div className="infoCard"><Info k="Телефон" v="Скрыт"/><Info k="О себе" v="Создатель AnonGram"/><Info k="AnonGram ID" v="#00000001"/><Info k="Звёзды" v="12 480"/><Info k="Участник с" v="28 сентября 2026"/></div><div className="musicCard"><div className="musicDisc">♪</div><div><b>AnonGram — Night Protocol</b><small>Музыка профиля</small></div><button>＋</button></div><h3>Публикации</h3><div className="gallery">{[1,2,3,4,5,6].map(x=><div key={x}><span>ANON</span></div>)}</div></section></main>
}
function Info({k,v}:{k:string;v:string}){return <div className="infoRow"><span>{k}</span><b>{v}</b></div>}

function ChatView({chat,back}:{chat:Chat;back:()=>void}) {
 const [text,setText]=useState(''); const [msgs,setMsgs]=useState([{me:false,t:'Добро пожаловать в AnonGram 👋'},{me:true,t:'Привет!'}]);
 return <main className="chatPage"><header className="chatHeader"><button className="roundBtn" onClick={back}>‹</button><Avatar c={chat}/><div><b>{chat.title}</b><small>{chat.online?'в сети':'был(а) недавно'}</small></div><div className="headRight"><button>⌕</button><button>⋮</button></div></header><div className="chatMessages">{msgs.map((m,i)=><div className={m.me?'bubble me':'bubble'} key={i}>{m.t}<small>12:{40+i}</small></div>)}</div><div className="composer"><button>＋</button><input value={text} onChange={e=>setText(e.target.value)} onKeyDown={e=>{if(e.key==='Enter'&&text.trim()){setMsgs([...msgs,{me:true,t:text.trim()}]);setText('')}}} placeholder="Сообщение"/><button onClick={()=>{if(text.trim()){setMsgs([...msgs,{me:true,t:text.trim()}]);setText('')}}}>➤</button></div></main>
}

function App(){
 const [logged,setLogged]=useState(false); const [drawer,setDrawer]=useState(false); const [page,setPage]=useState<'chats'|'profile'>('chats'); const [chat,setChat]=useState<Chat|null>(null); const [tab,setTab]=useState('Все чаты'); const [q,setQ]=useState('');
 const filtered=useMemo(()=>chats.filter(c=>(tab==='Все чаты'||c.type===tab|| (tab==='Группы'&&c.type==='Группы'))&&(!q||`${c.title} ${c.preview}`.toLowerCase().includes(q.toLowerCase()))),[tab,q]);
 if(!logged)return <Auth onLogin={()=>setLogged(true)}/>;
 if(page==='profile')return <Profile back={()=>setPage('chats')}/>;
 if(chat)return <ChatView chat={chat} back={()=>setChat(null)}/>;
 return <div className="app"><header className="topbar"><button className="roundBtn" onClick={()=>setDrawer(true)}>☰</button><Logo small/><div><button className="roundBtn">☾</button><button className="roundBtn">⋮</button></div></header>{drawer&&<Drawer close={()=>setDrawer(false)} onProfile={()=>{setDrawer(false);setPage('profile')}} logout={()=>{setDrawer(false);setLogged(false)}}/>}<main className="home"><section className="homeHero"><div className="homeHole"/><div className="heroText"><small>ANONGRAM · PRIVATE SPACE</small><h1>Познай<br/>анонимность_</h1><p>Общайся свободно. Создавай своё пространство.</p></div><img className="heroAvatar" src={avatar}/></section><div className="tabs">{['Все чаты','Группы','Каналы','Боты','Избранное'].map(x=><button className={tab===x?'active':''} onClick={()=>setTab(x)} key={x}>{x}</button>)}</div><div className="search"><span>⌕</span><input value={q} onChange={e=>setQ(e.target.value)} placeholder="Поиск по чатам"/></div><section className="chatList">{filtered.map(c=><article className="chatItem" key={c.id} onClick={()=>setChat(c)}><Avatar c={c}/><div className="chatMain"><div><b>{c.title}</b>{c.online&&<i className="verified">✓</i>}</div><p>{c.preview}</p></div><div className="chatMeta"><small>{c.time}</small>{c.unread>0&&<em>{c.unread}</em>}</div></article>)}</section></main><div className="createIsland"><div className="islandActions"><button>＋<span>Новый чат</span></button><button>♧<span>Группа</span></button><button>▣<span>Канал</span></button><button>◉<span>Анонимный</span></button></div><button className="mainPlus">＋</button></div></div>
}
createRoot(document.getElementById('root')!).render(<App/>);
