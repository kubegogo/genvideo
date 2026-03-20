import React from 'react'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import VideoGenerate from './pages/VideoGenerate'
import VideoRepurposing from './pages/VideoRepurposing'
import Settings from './pages/Settings'
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <nav className="sidebar">
          <h1>GenVideo</h1>
          <ul>
            <li><Link to="/">视频生成</Link></li>
            <li><Link to="/repurposing">视频搬运</Link></li>
            <li><Link to="/settings">设置</Link></li>
          </ul>
        </nav>
        <main className="content">
          <Routes>
            <Route path="/" element={<VideoGenerate />} />
            <Route path="/repurposing" element={<VideoRepurposing />} />
            <Route path="/settings" element={<Settings />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App
