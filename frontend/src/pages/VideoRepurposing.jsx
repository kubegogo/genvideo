import React, { useState } from 'react'
import axios from 'axios'

const API_BASE = '/api/v1'

function VideoRepurposing() {
  const [input, setInput] = useState('')
  const [inputType, setInputType] = useState('keywords')
  const [style, setStyle] = useState('documentary')
  const [duration, setDuration] = useState(60)
  const [aspectRatio, setAspectRatio] = useState('16:9')
  const [music, setMusic] = useState('relaxed')
  const [taskId, setTaskId] = useState(null)
  const [taskStatus, setTaskStatus] = useState(null)
  const [output, setOutput] = useState(null)

  // 生成视频
  const handleGenerate = async () => {
    try {
      const response = await axios.post(`${API_BASE}/video/generate`, {
        input,
        input_type: inputType,
        style,
        duration,
        aspect_ratio: aspectRatio,
        music,
      })
      const task = response.data.data
      setTaskId(task.id)
      setTaskStatus(task.status)
      setOutput(null)

      // 轮询任务状态
      const pollTask = async () => {
        try {
          const res = await axios.get(`${API_BASE}/task/${task.id}`)
          const updatedTask = res.data.data
          setTaskStatus(updatedTask.status)
          setOutput(updatedTask.output)

          if (updatedTask.status === 'completed') {
            return
          } else if (updatedTask.status === 'failed') {
            return
          }
          setTimeout(pollTask, 2000)
        } catch (err) {
          console.error('Poll error:', err)
        }
      }
      pollTask()
    } catch (error) {
      console.error('Generate failed:', error)
    }
  }

  // 下载视频
  const handleDownload = async () => {
    try {
      const response = await axios.post(`${API_BASE}/video/download`, {
        platform: 'douyin',
        video_url: input,
        metric_type: 'likes',
      })
      console.log('Download task:', response.data.data)
    } catch (error) {
      console.error('Download failed:', error)
    }
  }

  // 发布视频
  const handlePublish = async () => {
    if (!output) return
    try {
      await axios.post(`${API_BASE}/video/publish`, {
        video_path: `oss://genvideo/${output}`,
        platforms: ['youtube'],
        caption: 'AI Generated Video',
        tags: ['AI', 'generated'],
      })
    } catch (error) {
      console.error('Publish failed:', error)
    }
  }

  return (
    <div className="video-repurposing">
      <div className="card">
        <h2>AI 视频生成</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          输入文案或关键词，AI 自动生成视频素材并剪辑成片
        </p>

        <div className="form-group">
          <label>输入类型</label>
          <select value={inputType} onChange={(e) => setInputType(e.target.value)}>
            <option value="keywords">关键词</option>
            <option value="script">脚本/文案</option>
            <option value="article">文章</option>
          </select>
        </div>

        <div className="form-group">
          <label>输入内容</label>
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={
              inputType === 'keywords'
                ? '例如：未来城市、科幻、机器人'
                : '输入脚本或文案内容...'
            }
          />
        </div>

        <div className="form-group">
          <label>视频风格</label>
          <select value={style} onChange={(e) => setStyle(e.target.value)}>
            <option value="documentary">纪录片</option>
            <option value="dramatic">戏剧化</option>
            <option value="comedy">喜剧</option>
            <option value="action">动作</option>
            <option value="romance">浪漫</option>
          </select>
        </div>

        <div className="form-group">
          <label>时长（秒）</label>
          <input
            type="number"
            value={duration}
            onChange={(e) => setDuration(parseInt(e.target.value))}
            min={10}
            max={300}
          />
        </div>

        <div className="form-group">
          <label>画面比例</label>
          <select value={aspectRatio} onChange={(e) => setAspectRatio(e.target.value)}>
            <option value="16:9">横屏 (16:9)</option>
            <option value="9:16">竖屏 (9:16)</option>
            <option value="1:1">方形 (1:1)</option>
          </select>
        </div>

        <div className="form-group">
          <label>背景音乐风格</label>
          <select value={music} onChange={(e) => setMusic(e.target.value)}>
            <option value="relaxed">舒缓</option>
            <option value="energetic">活力</option>
            <option value="emotional">情感</option>
            <option value="cinematic">电影感</option>
            <option value="none">无音乐</option>
          </select>
        </div>

        <button className="btn btn-primary" onClick={handleGenerate}>
          生成视频
        </button>
      </div>

      {taskId && (
        <div className="card">
          <h2>任务 #{taskId}</h2>
          <div className="task-status-container">
            <span className={`task-status status-${taskStatus}`}>
              {taskStatus}
            </span>
          </div>

          {output && (
            <div style={{ marginTop: '20px' }}>
              <p>生成完成！视频路径: {output}</p>
              <button className="btn btn-primary" onClick={handlePublish}>
                发布视频
              </button>
            </div>
          )}
        </div>
      )}

      <div className="card">
        <h2>视频搬运</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          下载热门视频，进行二次创作后发布
        </p>

        <button className="btn btn-secondary" onClick={handleDownload}>
          下载视频
        </button>
      </div>
    </div>
  )
}

export default VideoRepurposing
