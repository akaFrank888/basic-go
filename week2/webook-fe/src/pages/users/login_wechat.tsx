import React, { useState, useEffect } from 'react';
import axios from "@/axios/axios";
import {redirect} from "next/navigation";

function Page() {
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        setLoading(true)
        axios.get('/oauth2/wechat/authurl')
            .then((res: { data: any; }) => res.data)
            .then((data: { data: string; }) => {
                setLoading(false)
                if(data && data.data) {
                    window.location.href = data.data
                }
            })
    }, [])

    if (isLoading) return <p>Loading...</p>

    return (
        <div>

        </div>
    )
}

export default Page