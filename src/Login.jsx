import { useEffect, useRef, useState } from "react"
import { Navigate, useNavigate } from "react-router-dom";


const Login = () => {
    const [error, setError] = useState(null);
    const user_id = useRef();
    const password = useRef();
    const navigate = useNavigate();

    useEffect(() => {
        const verify = async() => {
            try {
                const token = localStorage.getItem('token');
                const res = await fetch('http://localhost:8080/verify', {
                    headers: {'Authorization': `Bearer ${token}`}
                });
                if(res.ok){
                   navigate('/favorites');
                }
            } catch (err) {
                console.log(err);
                setError(err);
            }
        };
        verify();
    },[]);

    const handleLogin = async(e) => {
        try{
            e.preventDefault();
            const id = user_id.current.value;
            const pass = password.current.value;
            if(!id||!pass) {setError('ユーザーIDとパスワードを入力してください'); return;}
            const res = await fetch('http://localhost:8080/login', {
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
            navigate('/favorites');
        } catch(err) {
            console.log(err);
            setError(err);
        }
    }

    return(
        <div>
            <form onSubmit={handleLogin}>
                <div>
                    <label>
                        <span>User_ID</span>
                        <input type="text" ref={user_id}/>
                    </label>
                    <label>
                        <span>Password</span>
                        <input type="password" ref={password}/>
                    </label>
                </div>
                <div>
                    <button>Login</button>
                </div>
            </form>
            {error && <p>{error}</p>}
        </div>
    )

}

export default Login