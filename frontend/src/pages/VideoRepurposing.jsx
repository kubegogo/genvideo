import React, { useState } from 'react'
import axios from 'axios'

const API = '/api/v1'

export default function VideoRepurposing() {
  const [platform, setPlatform] = useState('douyin')
  const [videoUrl, setVideoUrl] = useState('')
  const [metricType, setMetricType] = useState('likes')

  const handleDownload = async () => {
    try {
      await axios.post(`${API}/video/download`, {
        platform,
        video_url: videoUrl,
        metric_type: metricType,
      })
      alert('下载任务已创建')
    } catch (err) {
      console.error('Download error:', err)
    }
  }

  return (
    <div>
      <div className="card">
        <h2>视频搬运</h2>
        <p style={{ color: '#666', marginBottom: 20 }}>
          从平台下载热门视频，进行AI二次创作后发布
        </p>

        <div className="form-group">
          <label>平台</label>
          <select value={platform} onChange={e => setPlatform(e.target.value)}>
            <option value="douyin">抖音</option>
            <option value="kuaishou">快手</option>
            <option value="bilibili">B站</option>
            <option value="xiaohongshu">小红书</option>
          </select>
        </div>

        <div className="form-group">
          <label>视频URL</label>
          <input type="text" value={videoUrl} onChange={e => setVideoUrl(e.target.value)} placeholder="输入视频链接" />
        </div>

        <div className="form-group">
          <label>热门指标</label>
          <select value={metricType} onChange={e => setMetricType(e.target.value)}>
            <option value="likes">点赞数</option>
            <option value="views">浏览数</option>
            <option value="favorites">收藏数</option>
          </select>
        </div>

        <button className="btn btn-primary" onClick={handleDownload}>
          下载视频
        </button>
      </div>
    </div>
  )
}
