import { BrowserRouter, Link, Route, Routes, Navigate } from 'react-router-dom'
import './App.css'
import Search from './Search';
import FavoriteList from './FavoriteList';
import Favorite from './Favorite'

function App() {
  return (
    <>
      {/*ルーティング機能を有効にする範囲を囲む　Link,Route*/}
      <BrowserRouter>
        <header>
          <h1>Video Search</h1>
          <nav className="header-nav">
            {/*ページリロードなしでURLを切り替えるリンク*/}
            <Link className='link' to={'/search'}>検索</Link>
            <Link className='link' to={'/favorites'}>お気に入り</Link>
          </nav>
        </header>
        <Routes>
          {/* リンクごとに表示する内容を切り替える */}
          {/* NavigateでURLを自動的に変更 */}
          <Route path='/' element={<Navigate to='/search'/>} />
          <Route path='/search' element={<Search />} />
          <Route path='/favorites' element={<FavoriteList/>} />
          <Route path='/favorites/:imdbID' element={<Favorite/>}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
