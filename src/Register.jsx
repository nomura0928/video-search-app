import { useRef, useState, useEffect } from "react"
import { Link, useNavigate } from "react-router-dom";
import './Auth.css';
import { verifyToken } from "./utils";


const Register = ({setIsLoggedIn}) => {
    const [error, setError] = useState(null);
    const user_id = useRef();
    const password = useRef();
    const navigate = useNavigate();

    useEffect(() => {
        const handleVerify = async () => {
            try {
                const ok = await verifyToken(setIsLoggedIn);
                if (ok) {
                    navigate('/favorites');
                }
            } catch (err) {
                console.log(err);
                setError(err);
            }
        };
        handleVerify();
    }, []);

    const handleRegister = async (e) => {
        try {
            e.preventDefault();
            const id = user_id.current.value;
            const pass = password.current.value;
            if (!id || !pass) { setError('ユーザーIDとパスワードを入力してください'); return; }
            const res = await fetch('https://video-search-app-3zcd.onrender.com/register', {
                headers: { 'Content-Type': 'application/json' },
                method: 'POST',
                body: JSON.stringify({
                    'user_id': id,
                    'password': pass
                })
            });
            if (!res.ok) {
                setError('登録に失敗しました');
                return;
            }
            navigate('/login');
        } catch (err) {
            console.log(err);
            setError(err);
        }
    }

    return (
        <div>
            <form onSubmit={handleRegister} className='auth-section'>
                <div className='input-left'>
                    <label className='search-div'>
                        <span className='input-info'>User ID</span>
                        <input type="text" ref={user_id} />
                    </label>
                    <label className='search-div'>
                        <span className='input-info'>Password</span>
                        <input type="password" ref={password} />
                    </label>
                </div>
                <button className='button'>Register</button>
            </form>
            <Link className='auth-link' to={'/login'}>ログインはこちら</Link>
            {error && <p className='auth-error'>{error}</p>}
        </div>
    )
}
export default Register
