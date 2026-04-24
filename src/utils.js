export const verifyToken = async (setIsLoggedIn) => {
    const token = localStorage.getItem('token');
    const res = await fetch('http://localhost:8080/verify', {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.status === 401) {
        setIsLoggedIn(false);
        return false;
    }
    setIsLoggedIn(true);
    return true;
};
