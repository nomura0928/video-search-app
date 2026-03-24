import React from 'react'

const ReviewList = ({ Ratings }) => {

    if (Ratings && Ratings.length > 0) {
        return (
            <>
                <table>
                    <tbody>
                        {Ratings.map((data) =>
                            <>
                            <tr>
                                <th className='source'>{data.Source} :</th>
                                <td>{data.Value}</td>
                            </tr>
                            </>
                        )}
                    </tbody>
                </table>
            </>
        )
    }
    else {
        return (
            <span>N/A</span>
        )
    }
}

export default ReviewList