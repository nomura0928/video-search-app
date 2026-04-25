import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import './FavoriteList.css'
import { verifyToken } from './utils';

const FavoriteList = ({ setIsLoggedIn }) => {
    const [favorites, setFavorites] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const getFavoriteList = async () => {
            try {
                const ok = await verifyToken(setIsLoggedIn);
                if (!ok) { navigate('/login'); return; }

                const token = localStorage.getItem('token');
                const res = await fetch('https://video-search-app-3zcd.onrender.com/favorites', {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                const data = await res.json();
                setFavorites(data);
            } catch (err) {
                console.log(err);
            } finally {
                //finallyは成功でも失敗でも必ず実行される
                setLoading(false);
            }
        };
        getFavoriteList();
    }, [])

    if (loading) return null;

    if (favorites.length === 0) {
        return <h2>お気に入り登録中の映画はありません</h2>
    } else {
        return (
            <div className='favorites-grid'>
                {favorites.map((favorite) =>
                    <Link key={favorite.imdbID} className='favorite-card' to={`/favorites/${favorite.imdbID}`}>
                        <img src={favorite.Poster} alt="Poster" onError={(e) => e.target.src = '/src/assets/noimage.jpg'} />
                        <h4>{favorite.Title}</h4>
                        <h4>{favorite.Year}</h4>
                    </Link>
                )}
            </div>
        )
    }
}

export default FavoriteList
