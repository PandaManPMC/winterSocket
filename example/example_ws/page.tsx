"use client";

import React, { useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";

const HomePage = () => {
    const [wsUrl, setWsUrl] = useState("ws://localhost:19999");
    const [connectUrl, setConnectUrl] = useState("");
    const { sendMessage, lastMessage, readyState } = useWebSocket(connectUrl, {
        shouldReconnect: () => true,
        reconnectInterval: 3000,
    }, connectUrl !== "");

    const [latestMessage, setLatestMessage] = useState<string | null>(null);
    const [inputMessage, setInputMessage] = useState("login&"+JSON.stringify({ userToken: "abc中文" }));

    useEffect(() => {
        console.log("WebSocket 状态变化:", readyState);
    }, [readyState]);

    useEffect(() => {
        if (!lastMessage) return;
        console.log("收到消息:", lastMessage.data);

        if (lastMessage.data instanceof Blob) {
            lastMessage.data.text().then(setLatestMessage);
        } else {
            setLatestMessage(lastMessage.data);
        }
    }, [lastMessage]);

    const handleSendMessage = () => {
        if (inputMessage.trim() !== "") {
            sendMessage(inputMessage);
            setInputMessage("");
        }
    };

    const handleConnect = () => {
        setConnectUrl(wsUrl);
    };

    return (
        <div>
            <input
                type="text"
                value={wsUrl}
                onChange={(e) => setWsUrl(e.target.value)}
                placeholder="输入 WebSocket 连接地址..."
                style={{ width: "100%", margin: "10px", borderWidth: "1px" }}
            />
            <button style={{ backgroundColor: "#90EE90", margin: "10px" }} onClick={handleConnect}>连接 WebSocket</button>
            <input
                type="text"
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                placeholder="输入消息..."
                style={{ width: "100%", margin: "10px", borderWidth: "1px" }}
            />
            <div>
                <button style={{ backgroundColor: "#ADD8E6", margin: "10px" }} onClick={handleSendMessage}>发送消息</button>
            </div>
            <p>最新消息: {latestMessage}</p>
            <p>
                连接状态: {" "}
                {readyState === WebSocket.OPEN
                    ? "✅ 连接成功"
                    : readyState === WebSocket.CONNECTING
                        ? "⏳ 连接中..."
                        : "❌ 连接关闭"}
            </p>
            <div style={{ margin: "10px" }}>
                <hr/>
                <p>ping</p>
                <hr/>
                <p>login&{JSON.stringify({ userToken: "abc中文" })}</p>
                <hr/>
            </div>
        </div>
    );
};

export default HomePage;