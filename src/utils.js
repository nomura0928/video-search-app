export const verifyToken = async (setIsLoggedIn) => {
    const token = localStorage.getItem('token');
    const res = await fetch('https://video-search-app-3zcd.onrender.com/verify', {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.status === 401) {
        setIsLoggedIn(false);
        return false;
    }
    setIsLoggedIn(true);
    return true;
};
