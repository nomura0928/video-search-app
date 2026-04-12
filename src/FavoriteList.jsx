import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import './FavoriteList.css'

const FavoriteList = () => {
  const [favorites,setFavorites] = useState([]);

  useEffect(() => {
    try {
      fetch(`http://localhost:8080/favorites`)
      .then(res => res.json())
      .then(data => setFavorites(data))
    } catch (err) {
      console.log(err);
    }
  },[])

  return (
    <div className='favorites-grid'>
      {favorites.map((favorite) =>
        <Link key={favorite.imdbID} className='favorite-card' to={`/favorites/${favorite.imdbID}`}>
          <img src={favorite.Poster} alt="Poster" onError={(e) => e.target.src = '/src/assets/noimage.jpg'}/>
          <h4>{favorite.Title}</h4>
          <h4>{favorite.Year}</h4>
        </Link>
      )}
    </div>
  )
}

export default FavoriteList