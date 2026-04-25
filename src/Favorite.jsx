import { useEffect, useState } from 'react';
import FilmDetail from './FilmDetail';
import { useNavigate, useParams } from 'react-router-dom';
import { verifyToken } from './utils';

const Favorite = ({ isLoggedIn, setIsLoggedIn }) => {
    const [filmInfo, setFilmInfo] = useState();
    const { imdbID } = useParams();
    const navigate = useNavigate();

    useEffect(() => {
        const HandleGetInfo = async () => {
            try {
                const ok = await verifyToken(setIsLoggedIn);
                if (!ok) {navigate('/login'); return;}
                const res = await fetch('https://video-search-app-3zcd.onrender.com/search', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 'imdbID': imdbID })
                });
                const data = await res.json();
                setFilmInfo(data);
            } catch (err) {
                console.log(err);
            }
        };
        HandleGetInfo();
    }, [])

    if (!filmInfo) {
        return <p>読み込み中</p>
    } else if (filmInfo.Response) {
        return <FilmDetail data={filmInfo} isLoggedIn={isLoggedIn} setIsLoggedIn={setIsLoggedIn} />
    } else {
        return <h2>取得に失敗しました</h2>
    }
}

export default Favorite
