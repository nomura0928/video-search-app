import React, { useEffect } from 'react';
import ReviewList from './ReviewList';

const Information = ({ data }) => {

    // const imgcheck = async (url) => {
    //     try{
    //         const res = await fetch(url);
    //     } catch (err){
    //         data.Poster = "src/assets/noimage.jpg";
    //     }
    // }

    // useEffect(() => {
    //     if(data.Poster === 'N/A'){
    //         data.Poster = "src/assets/noimage.jpg";
    //     }
    //     else {
    //         imgcheck(data.Poster);
    //     }
    // }, [data])

    if (data && data.Response === 'True') {
        return (
            <>
                <figure className='info-top'>
                    <div className='poster-div'>
                    {/* posterが表示できないとき、代替画像を表示 */}
                        <img src={data.Poster} alt="poster" onError={(e) => e.target.src = '/src/assets/noimage.jpg'} />
                        <button className='favorite-button'>お気に入り登録</button>
                    </div>
                    <figcaption className='figcaption'>
                        <table>
                            <tbody>
                                <tr className='tr'>
                                    <th className='th'>タイトル</th>
                                    <td className='td'>{data.Title}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>公開年</th>
                                    <td className='td'>{data.Year}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>ジャンル</th>
                                    <td className='td'>{data.Genre}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>形式</th>
                                    <td className='td'>{data.Type}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>監督</th>
                                    <td className='td'>{data.Director}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>脚本</th>
                                    <td className='td'>{data.Writer}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th'>出演</th>
                                    <td className='td'>{data.Actors}</td>
                                </tr>
                                <tr className='tr'>
                                    <th className='review th'>評価</th>
                                    <td className='td'>
                                            <ReviewList Ratings={data.Ratings} />
                                    </td>
                                </tr>
                                <tr className='tr'>
                                    <th className='th review'>概要</th>
                                    <td className='td'>{data.Plot}</td>
                                </tr>
                            </tbody>
                        </table>
                    </figcaption>
                </figure>
            </>
        )
    }
    else if(!data){
        return(
            <>
            <h2>探したい映画やドラマの名前や公開年を英語で入力してください</h2>
            <h3>少なくともタイトルは入力してください</h3>
            </>
        )
    }
    else {
        return (
            <h2>そのような映像作品はデータベース上にありません</h2>
        )
    }
}

export default Information