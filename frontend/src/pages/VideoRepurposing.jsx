import React, { useState } from 'react'

const API_BASE = '/api'

function VideoRepurposing() {
  const [platform, setPlatform] = useState('douyin')
  const [videoUrl, setVideoUrl] = useState('')
  const [style, setStyle] = useState('cartoon')
  const [loading, setLoading] = useState(false)
  const [taskId, setTaskId] = useState(null)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const handleDownload = async () => {
    if (!videoUrl) {
      setError('请输入视频链接')
      return
    }
    setLoading(true)
    setError(null)
    try {
      const res = await fetch(`${API_BASE}/download`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ platform, video_url: videoUrl, metric_type: 'likes' })
      })
      const data = await res.json()
      if (data.code === 0) {
        setTaskId(data.data.task_id)
        pollTaskStatus(data.data.task_id)
      } else {
        setError(data.msg)
      }
    } catch (err) {
      setError('请求失败: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const pollTaskStatus = async (id) => {
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${API_BASE}/task/${id}`)
        const data = await res.json()
        if (data.code === 0) {
          const task = data.data
          if (task.status === 'completed') {
            setResult(task.output)
            clearInterval(interval)
          } else if (task.status === 'failed') {
            setError(task.error)
            clearInterval(interval)
          }
        }
      } catch (err) {
        console.error('Poll error:', err)
      }
    }, 3000)
  }

  return (
    <div className="page">
      <h1>视频搬运</h1>
      <p className="subtitle">下载热门视频 → AI二次创作 → 发布到其他平台</p>

      <div className="card">
        <div className="form-group">
          <label>选择平台</label>
          <select value={platform} onChange={e => setPlatform(e.target.value)}>
            <option value="douyin">抖音</option>
            <option value="kuaishou">快手</option>
            <option value="bilibili">B站</option>
            <option value="xiaohongshu">小红书</option>
          </select>
        </div>

        <div className="form-group">
          <label>视频链接</label>
          <input
            type="text"
            placeholder="粘贴视频链接"
            value={videoUrl}
            onChange={e => setVideoUrl(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label>创作风格</label>
          <select value={style} onChange={e => setStyle(e.target.value)}>
            <option value="cartoon">卡通风格</option>
            <option value="anime">动漫风格</option>
            <option value="realistic">写实风格</option>
            <option value="abstract">抽象风格</option>
          </select>
        </div>

        <button className="btn-primary" onClick={handleDownload} disabled={loading}>
          {loading ? '处理中...' : '开始下载并创作'}
        </button>

        {error && <div className="error-msg">{error}</div>}

        {result && (
          <div className="result">
            <h3>创作完成!</h3>
            <a href={result} target="_blank" rel="noopener noreferrer">下载视频</a>
          </div>
        )}
      </div>
    </div>
  )
}

export default VideoRepurposing
