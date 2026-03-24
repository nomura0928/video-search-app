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

  const getInfo = async () => {
    const name = filmName.current.value;
    const year = filmYear.current.value;
    if (!name) return;
    try {
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
      <BrowserRouter>
        <header>
          <h1>Video Search</h1>
          <nav className="header-nav">
            <Link className='link' to={'/'}>検索</Link>
            <Link className='link' to={'/favorite'}>お気に入り</Link>
          </nav>
        </header>
        <section className='search-section'>
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
            <button className='button' onClick={getInfo}>Search</button>
          </div>
        </section>
        <section className='film-info'>
          <Routes>
          <Route path='/' element={<Information data={filmInfo} />} />
          <Route path='/favorite' element={<FavoriteList/>} />
          </Routes>
        </section>
      </BrowserRouter>
    </>
  )
}

export default App
