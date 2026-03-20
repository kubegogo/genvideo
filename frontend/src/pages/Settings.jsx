import React, { useState, useEffect } from 'react'
import axios from 'axios'

const API_BASE = '/api/v1'

function Settings() {
  const [aiProvider, setAiProvider] = useState('minimax')
  const [minimaxKey, setMinimaxKey] = useState('')
  const [n8nUrl, setN8nUrl] = useState('http://localhost:5678')
  const [comfyuiUrl, setComfyuiUrl] = useState('http://localhost:8188')
  const [ollamaUrl, setOllamaUrl] = useState('http://localhost:11434')
  const [ossEndpoint, setOssEndpoint] = useState('oss-cn-hangzhou.aliyuncs.com')
  const [ossBucket, setOssBucket] = useState('genvideo')
  const [ossKey, setOssKey] = useState('')
  const [ossSecret, setOssSecret] = useState('')
  const [videoProviders, setVideoProviders] = useState([])

  useEffect(() => {
    // Load current config
    const loadConfig = async () => {
      try {
        const providersRes = await axios.get(`${API_BASE}/config/ai-providers`)
        console.log('AI Providers:', providersRes.data)
      } catch (error) {
        console.error('Failed to load config:', error)
      }
    }
    loadConfig()
  }, [])

  const handleSaveAIProvider = async () => {
    try {
      await axios.post(`${API_BASE}/config/ai-providers`, {
        type: aiProvider,
        api_key: minimaxKey,
        base_url: aiProvider === 'self_hosted' ? n8nUrl : '',
        is_active: true,
      })
      alert('AI 配置已保存')
    } catch (error) {
      console.error('Failed to save AI config:', error)
    }
  }

  const handleSaveOSS = async () => {
    try {
      await axios.post(`${API_BASE}/config/oss`, {
        endpoint: ossEndpoint,
        access_key: ossKey,
        secret_key: ossSecret,
        bucket: ossBucket,
        is_active: true,
      })
      alert('OSS 配置已保存')
    } catch (error) {
      console.error('Failed to save OSS config:', error)
    }
  }

  return (
    <div className="settings">
      <div className="card">
        <h2>AI 服务配置</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          选择视频生成方式
        </p>

        <div className="form-group">
          <label>AI 提供商</label>
          <select value={aiProvider} onChange={(e) => setAiProvider(e.target.value)}>
            <option value="minimax">Minimax API</option>
            <option value="self_hosted">自建 (n8n + ComfyUI + Ollama)</option>
          </select>
        </div>

        {aiProvider === 'minimax' && (
          <div className="form-group">
            <label>Minimax API Key</label>
            <input
              type="password"
              value={minimaxKey}
              onChange={(e) => setMinimaxKey(e.target.value)}
              placeholder="输入您的 API Key"
            />
          </div>
        )}

        {aiProvider === 'self_hosted' && (
          <>
            <div className="form-group">
              <label>n8n 地址</label>
              <input
                type="text"
                value={n8nUrl}
                onChange={(e) => setN8nUrl(e.target.value)}
              />
            </div>
            <div className="form-group">
              <label>ComfyUI 地址</label>
              <input
                type="text"
                value={comfyuiUrl}
                onChange={(e) => setComfyuiUrl(e.target.value)}
              />
            </div>
            <div className="form-group">
              <label>Ollama 地址</label>
              <input
                type="text"
                value={ollamaUrl}
                onChange={(e) => setOllamaUrl(e.target.value)}
              />
            </div>
          </>
        )}

        <button className="btn btn-primary" onClick={handleSaveAIProvider}>
          保存 AI 配置
        </button>
      </div>

      <div className="card">
        <h2>阿里云 OSS 配置</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          视频文件将同步到阿里云 OSS
        </p>

        <div className="form-group">
          <label>OSS Endpoint</label>
          <input
            type="text"
            value={ossEndpoint}
            onChange={(e) => setOssEndpoint(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label>Bucket</label>
          <input
            type="text"
            value={ossBucket}
            onChange={(e) => setOssBucket(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label>Access Key</label>
          <input
            type="password"
            value={ossKey}
            onChange={(e) => setOssKey(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label>Secret Key</label>
          <input
            type="password"
            value={ossSecret}
            onChange={(e) => setOssSecret(e.target.value)}
          />
        </div>

        <button className="btn btn-primary" onClick={handleSaveOSS}>
          保存 OSS 配置
        </button>
      </div>

      <div className="card">
        <h2>视频平台配置</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          配置各平台的登录信息
        </p>

        {['抖音', '快手', 'B站', '小红书'].map((platform) => (
          <div key={platform} className="form-group">
            <label>{platform} Cookie</label>
            <input
              type="password"
              placeholder={`输入 ${platform} 的登录 Cookie`}
            />
          </div>
        ))}

        <button className="btn btn-primary">
          保存平台配置
        </button>
      </div>
    </div>
  )
}

export default Settings
