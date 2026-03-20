import React, { useState } from 'react'
import axios from 'axios'

const API_BASE = '/api/v1'

function VideoRepurposing() {
  const [platform, setPlatform] = useState('douyin')
  const [videoUrl, setVideoUrl] = useState('')
  const [metricType, setMetricType] = useState('likes')
  const [style, setStyle] = useState('original')
  const [tasks, setTasks] = useState([])

  const handleDownload = async () => {
    try {
      const response = await axios.post(`${API_BASE}/video/download`, {
        platform,
        video_url: videoUrl,
        metric_type: metricType,
      })
      setTasks([...tasks, response.data.data])
    } catch (error) {
      console.error('Download failed:', error)
    }
  }

  const handleRecreate = async () => {
    try {
      const response = await axios.post(`${API_BASE}/video/recreate`, {
        original_video: videoUrl,
        style,
        keep_audio: true,
      })
      setTasks([...tasks, response.data.data])
    } catch (error) {
      console.error('Recreate failed:', error)
    }
  }

  const handlePublish = async (taskId) => {
    try {
      const response = await axios.post(`${API_BASE}/video/publish`, {
        video_path: `oss://genvideo/output/${taskId}/video.mp4`,
        platforms: ['youtube'],
        caption: 'AI generated video',
        tags: ['AI', 'generated'],
      })
      setTasks([...tasks, response.data.data])
    } catch (error) {
      console.error('Publish failed:', error)
    }
  }

  return (
    <div className="video-repurposing">
      <div className="card">
        <h2>视频搬运</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          从平台下载热门视频，AI 二次创作后发布到其他平台
        </p>

        <div className="form-group">
          <label>选择平台</label>
          <select value={platform} onChange={(e) => setPlatform(e.target.value)}>
            <option value="douyin">抖音</option>
            <option value="kuaishou">快手</option>
            <option value="bilibili">B站</option>
            <option value="xiaohongshu">小红书</option>
          </select>
        </div>

        <div className="form-group">
          <label>视频 URL 或搜索条件</label>
          <input
            type="text"
            value={videoUrl}
            onChange={(e) => setVideoUrl(e.target.value)}
            placeholder="输入视频链接或关键词"
          />
        </div>

        <div className="form-group">
          <label>选择热门指标</label>
          <select value={metricType} onChange={(e) => setMetricType(e.target.value)}>
            <option value="likes">点赞数</option>
            <option value="views">浏览数</option>
            <option value="favorites">收藏数</option>
          </select>
        </div>

        <button className="btn btn-primary" onClick={handleDownload}>
          下载视频
        </button>
      </div>

      <div className="card">
        <h2>二次创作</h2>

        <div className="form-group">
          <label>创作风格</label>
          <select value={style} onChange={(e) => setStyle(e.target.value)}>
            <option value="original">保持原风格</option>
            <option value="dramatic">戏剧化</option>
            <option value="comedy">喜剧化</option>
            <option value="documentary">纪录片风格</option>
          </select>
        </div>

        <button className="btn btn-primary" onClick={handleRecreate}>
          AI 二次创作
        </button>
      </div>

      {tasks.length > 0 && (
        <div className="card">
          <h2>任务列表</h2>
          <div className="task-list">
            {tasks.map((task, index) => (
              <div key={index} className="task-item">
                <div>
                  <strong>任务 #{task.id}</strong>
                  <p style={{ fontSize: '12px', color: '#666' }}>{task.type}</p>
                </div>
                <div>
                  <span className={`task-status status-${task.status}`}>
                    {task.status}
                  </span>
                  {task.status === 'completed' && task.type === 'repurposing' && (
                    <button
                      className="btn btn-secondary"
                      style={{ marginLeft: '10px' }}
                      onClick={() => handlePublish(task.id)}
                    >
                      发布
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default VideoRepurposing
