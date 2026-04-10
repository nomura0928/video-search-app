import { useRef, useState } from 'react';
import ReviewList from './ReviewList';
import './Information.css';

const Information = () => {
    const apikey = import.meta.env.VITE_OMDB_API_KEY;
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
            {filmInfo && filmInfo.Response === 'True' ? (
                <figure className='info-top'>
                    <div className='poster-div'>
                        {/* posterが表示できないとき、代替画像を表示 */}
                        <img src={filmInfo.Poster} alt="poster" onError={(e) => e.target.src = '/src/assets/noimage.jpg'} />
                        <button className='favorite-button'>お気に入り登録</button>
                    </div>
                    <figcaption className='figcaption'>
                        <table>
                            <tbody>
                                <tr className='tr'>
                                    <th className='th'>タイトル</th>
                                    <td className='td'>{filmInfo.Title}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>公開年</th>
                                    <td className='td'>{filmInfo.Year}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>ジャンル</th>
                                    <td className='td'>{filmInfo.Genre}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>形式</th>
                                    <td className='td'>{filmInfo.Type}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>監督</th>
                                    <td className='td'>{filmInfo.Director}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>脚本</th>
                                    <td className='td'>{filmInfo.Writer}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>出演</th>
                                    <td className='td'>{filmInfo.Actors}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='review th'>評価</th>
                                    <td className='td'>
                                        <ReviewList Ratings={filmInfo.Ratings} />
                                    </td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th review'>概要</th>
                                    <td className='td'>{filmInfo.Plot}</td>
                                </tr>
                            </tbody>
                        </table>
                    </figcaption>
                </figure>
            ) : !filmInfo ? (
                <>
                    <h2>探したい映画やドラマの名前や公開年を英語で入力してください</h2>
                    <h3>少なくともタイトルは入力してください</h3>
                </>
            ) : (
                <h2>そのような映像作品はデータベース上にありません</h2>
            )}
        </>
    )
}

export default Information
