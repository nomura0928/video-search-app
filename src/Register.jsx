import { useRef, useState, useEffect } from "react"
import { useNavigate } from "react-router-dom";


const Register = () => {
    const [error, setError] = useState(null);
    const user_id = useRef();
    const password = useRef();
    const navigate = useNavigate();

    useEffect(() => {
            const handleVerify = async() => {
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
            handleVerify();
        },[]);

        const handleRegister = async(e) => {
            try{
                e.preventDefault();
                const id = user_id.current.value;
                const pass = password.current.value;
                if(!id||!pass) {setError('ユーザーIDとパスワードを入力してください'); return;}
                const res = await fetch('http://localhost:8080/register', {
                    headers: {'Content-Type': 'application/json'},
                    method: 'POST',
                    body: JSON.stringify({
                        'user_id': id,
                        'password': pass
                    })
                });
                if(!res.ok){
                    setError('登録に失敗しました');
                    return;
                }
                const data = await res.json();
                navigate('/login');
            } catch(err) {
                console.log(err);
                setError(err);
            }
        }

    return(
        <div>
            <form onSubmit={handleRegister}>
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
                    <button>Register</button>
                </div>
            </form>
            {error && <p>{error}</p>}
        </div>
    )
}
export default Register