import { useEffect, useRef, useState } from 'react';
import FilmDetail from './FilmDetail';
import './Search.css';

const Search = ({ isLoggedIn, setIsLoggedIn }) => {
    const filmName = useRef();
    const filmYear = useRef();
    const [filmInfo, setInfo] = useState(null);
    useEffect(() => {
        console.log(filmInfo);
    }, [filmInfo])

    const getInfo = async (e) => {
        //submitのリロードを無効化
        e.preventDefault();
        const name = filmName.current.value;
        const year = filmYear.current.value;
        if (!name) return;
        try {
            //awaitで非同期処理が完了するまで待つ
            const res = await fetch(`https://video-search-app-3zcd.onrender.com/search`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({'Title': name,'Year': year})
            });
            const data = await res.json();
            setInfo(data);
        } catch (err) {
            console.log(err);
        }
    }

    return (
        <>
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
            <FilmDetail data={filmInfo} isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn} />
        </>
    )
}

export default Search
