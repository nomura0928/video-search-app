import { useEffect, useState } from 'react';
import FilmDetail from './FilmDetail';
import { useParams } from 'react-router-dom';

const Favorite = () => {
    //Hookはトップレベルで呼ぶ必要がある
    const [filmInfo,setFilmInfo] = useState();
    const {imdbID} = useParams();

    useEffect(() => {
        const HandleGetInfo = async() => {
            try{
        const res = await fetch(`http://localhost:8080/search`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({'imdbID': imdbID})
        });
        const data = await res.json();
        setFilmInfo(data);
    } catch(err) {
        console.log(err);
    }
        }
        HandleGetInfo();
    },[])

    if (!filmInfo) {
        return <p>読み込み中</p>
    } else if(filmInfo.Response){
        return(
        <>
        <FilmDetail data={filmInfo}/>
        </>
        )
    } else {
        return <h2>取得に失敗しました</h2>
    }
}

export default Favorite