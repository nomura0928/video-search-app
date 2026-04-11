import React, { useEffect, useState } from 'react'

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
    <div>
      {favorites.map((favorites) =>
      <div>
        <img src={favorites.Poster} alt="Poster" onError={(e) => e.target.src = '/src/assets/noimage.jpg'}/>
        <h4>{favorites.Title}</h4>
        <h4>{favorites.Year}</h4>
      </div>
      )}
    </div>
  )
}

export default FavoriteList