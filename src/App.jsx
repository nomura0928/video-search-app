import { BrowserRouter, Link, Route, Routes, Navigate, useNavigate, useLocation } from 'react-router-dom'
import { useState } from 'react'
import './App.css'
import Search from './Search';
import FavoriteList from './FavoriteList';
import Favorite from './Favorite'
import Login from './Login';
import Register from './Register';

function Header({ isLoggedIn, setIsLoggedIn }) {
  const navigate = useNavigate();
  //現在のパスを取得
  const location = useLocation();

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    if(location.pathname !== '/search') navigate('/login');
  };

  return (
    <header>
      <h1>Video Search</h1>
      <nav className="header-nav">
        <Link className='link' to={'/search'}>検索</Link>
        <Link className='link' to={'/favorites'}>お気に入り</Link>
      </nav>
      <div className='header-auth'>
        {isLoggedIn ? (
          <button className='header-auth-btn' onClick={handleLogout}>ログアウト</button>
        ) : (
          <Link className='link' to={'/login'}>ログイン</Link>
        )}
      </div>
    </header>
  );
}

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));

  return (
    <>
      {/*ルーティング機能を有効にする範囲を囲む　Link,Route*/}
      <BrowserRouter>
        <Header isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn} />
        <Routes>
          {/* リンクごとに表示する内容を切り替える */}
          {/* NavigateでURLを自動的に変更 */}
          <Route path='/' element={<Navigate to='/search'/>} />
          <Route path='/search' element={<Search isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn}/>} />
          <Route path='/favorites' element={<FavoriteList setIsLoggedIn={setIsLoggedIn}/>} />
          <Route path='/favorites/:imdbID' element={<Favorite isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn}/>}/>
          <Route path='/login' element={<Login setIsLoggedIn={setIsLoggedIn}/>}/>
          <Route path='/register' element={<Register setIsLoggedIn={setIsLoggedIn}/>}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
