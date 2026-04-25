import { useEffect, useRef, useState } from "react"
import { useNavigate, Link } from "react-router-dom";
import './Auth.css';
import { verifyToken } from "./utils";


const Login = ({setIsLoggedIn}) => {
    const [error, setError] = useState(null);
    const user_id = useRef();
    const password = useRef();
    const navigate = useNavigate();

    useEffect(() => {
        const handleVerify = async() => {
            try {
                const ok = await verifyToken(setIsLoggedIn);
                if(ok){
                    navigate('/favorites');
                }
            } catch (err) {
                console.log(err);
                setError(err);
            }
        };
        handleVerify();
    },[]);

    const handleLogin = async(e) => {
        try{
            e.preventDefault();
            const id = user_id.current.value;
            const pass = password.current.value;
            if(!id||!pass) {setError('ユーザーIDとパスワードを入力してください'); return;}
            const res = await fetch('https://video-search-app-3zcd.onrender.com/login', {
                headers: {'Content-Type': 'application/json'},
                method: 'POST',
                body: JSON.stringify({
                    'user_id': id,
                    'password': pass
                })
            });
            if(!res.ok){
                setError('ユーザーIDまたはパスワードが違います');
                return;
            }
            const data = await res.json();
            localStorage.setItem('token', data.token);
            setIsLoggedIn(true);
            navigate('/favorites');
        } catch(err) {
            console.log(err);
            setError(err);
        }
    }

    return(
        <div>
            <form onSubmit={handleLogin} className='auth-section'>
                <div className='input-left'>
                    <label className='search-div'>
                        <span className='input-info'>User ID</span>
                        <input type="text" ref={user_id}/>
                    </label>
                    <label className='search-div'>
                        <span className='input-info'>Password</span>
                        <input type="password" ref={password}/>
                    </label>
                </div>
                <button className='button'>Login</button>
            </form>
            <Link className='auth-link' to={'/register'}>未登録の方はこちら</Link>
            {error && <p className='auth-error'>{error}</p>}
        </div>
    )

}

export default Login