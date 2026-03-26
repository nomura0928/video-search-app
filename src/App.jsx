import { useEffect, useRef, useState } from 'react'
import { BrowserRouter, Link, Route, Routes } from 'react-router-dom'
import './App.css'
import Information from './Information';
import FavoriteList from './FavoriteList';

function App() {
  const apikey = import.meta.env.VITE_OMDB_API_KEY
  const filmName = useRef();
  const filmYear = useRef();
  const [filmInfo, setInfo] = useState(null);

  const getInfo = async (e) => {
    e.preventDefault();
    const name = filmName.current.value;
    const year = filmYear.current.value;
    if (!name) return;
    try {
      //awaitで非同期処理が完了するまで待つ
      const res = await fetch(`https://www.omdbapi.com/?apikey=${apikey}&t=${name}&y=${year}`);
      console.log(res);
      const data = await res.json();
      setInfo(data);
    } catch (err) {
      console.log(err);
    }
  }

  useEffect(() => console.log(filmInfo), [filmInfo]);

  return (
    <>
    {/*ルーティング機能を有効にする範囲を囲む　Link,Route*/}
      <BrowserRouter>
        <header>
          <h1>Video Search</h1>
          <nav className="header-nav">
            {/*ページリロードなしでURLを切り替えるリンク*/}
            <Link className='link' to={'/'}>検索</Link>
            <Link className='link' to={'/favorite'}>お気に入り</Link>
          </nav>
        </header>
        <form onSubmit={getInfo} className='search-section'>
          <div className="input-left">
            <label className='search-div'>
              <span className='input-info'>タイトル</span>
              <input type="text" ref={filmName} placeholder='英語で入力' />
            </label>
            <label className='search-div'>
              <span className='input-info'>公開年</span>
              <input type="text" ref={filmYear} />
            </label>
          </div>
          <div className="input-right">
            <button className='button' type='submit'>Search</button>
          </div>
        </form>
        <section className='film-info'>
          <Routes>
            {/* リンクごとに表示する内容を切り替える */}
          <Route path='/' element={<Information data={filmInfo} />} />
          <Route path='/favorite' element={<FavoriteList/>} />
          </Routes>
        </section>
      </BrowserRouter>
    </>
  )
}

export default App
