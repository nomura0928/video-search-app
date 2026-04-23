export const verifyToken = async () => {
    const token = localStorage.getItem('token');
    const res = await fetch('http://localhost:8080/verify', {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.status === 401) {
        return false;
    }
    return true;
};
