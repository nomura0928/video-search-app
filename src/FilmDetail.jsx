import { useEffect, useState } from 'react';
import ReviewList from './ReviewList';

const FilmDetail = ({ data, isLoggedIn, setIsLoggedIn }) => {

    const [isFavorite,setIsFavorite] = useState(false);
    const token = localStorage.getItem('token');

    useEffect(() => {
        if(!data) return;
        if(!isLoggedIn) return;
        //useEffectに直接asyncはつけられない
        const check = async() => {
        try{
            const res = await fetch(`${import.meta.env.VITE_API_URL}/favorites`, {
                headers: {'Authorization': `Bearer ${token}`}
            });
            const favorites = await res.json();
            //some: 配列の中に条件を満たすものがあればtrueを返す
            setIsFavorite(favorites.some(f => f.imdbID === data.imdbID));
        } catch(err) {
            console.log(err);
        }
    }
        check();
    },[data, isLoggedIn]);

    const FavoriteButton = async() => {
        if(isFavorite){
            try{
                await fetch(`${import.meta.env.VITE_API_URL}/favorites`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(data)
                });
                setIsFavorite(false);
            } catch(err) {
                console.log(err);
            }
        } else {
            try{
                await fetch(`${import.meta.env.VITE_API_URL}/favorites`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(data)
                });
                setIsFavorite(true);
            } catch(err) {
                console.log(err);
            }
        }
    }

    if (data && data.Response === 'True') {
        return (
            <figure className='info-top'>
                <div className='poster-div'>
                    {/* posterが表示できないとき、代替画像を表示 */}
                    <img src={data.Poster} alt="poster" onError={(e) => e.target.src = '/noimage.jpg'} />
                    {/* onclick={関数()}とするとレンダリング後即実行されてしまう、()はつけない */}
                    {isLoggedIn && <button className='favorite-button' onClick={FavoriteButton}>{isFavorite ? 'お気に入り解除' : 'お気に入り登録'}</button>}
                </div>
                <figcaption className='figcaption'>
                    <table>
                        <tbody>
                            <tr className='tr'>
                                <th className='th'>タイトル</th>
                                <td className='td'>{data.Title}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>公開年</th>
                                <td className='td'>{data.Year}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>ジャンル</th>
                                <td className='td'>{data.Genre}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>形式</th>
                                <td className='td'>{data.Type}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>監督</th>
                                <td className='td'>{data.Director}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>脚本</th>
                                <td className='td'>{data.Writer}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='th'>出演</th>
                                <td className='td'>{data.Actors}</td>
                            </tr>
                            <tr className='tr'>
                                <th className='review th'>評価</th>
                                <td className='td'>
                                    <ReviewList Ratings={data.Ratings} />
                                </td>
                            </tr>
                            <tr className='tr'>
                                <th className='th review'>概要</th>
                                <td className='td'>{data.Plot}</td>
                            </tr>
                        </tbody>
                    </table>
                </figcaption>
            </figure>
        )
    }
    else if (!data) {
        return (
            <>
                <h2>探したい映画やドラマの名前や公開年を英語で入力してください</h2>
                <h3>少なくともタイトルは入力してください</h3>
            </>
        )
    }
    else {
        return (
            <h2>そのような映像作品はデータベース上にありません</h2>
        )
    }
}

export default FilmDetail
