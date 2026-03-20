import React from 'react'

export default function Settings() {
  return (
    <div>
      <div className="card">
        <h2>AI服务配置</h2>
        <p style={{ color: '#666', marginBottom: 20 }}>
          选择视频生成方式
        </p>

        <div className="form-group">
          <label>AI提供商</label>
          <select>
            <option value="minimax">Minimax API</option>
            <option value="self_hosted">自建 (n8n+ComfyUI+Ollama)</option>
          </select>
        </div>

        <div className="form-group">
          <label>API Key</label>
          <input type="password" placeholder="输入API Key" />
        </div>

        <button className="btn btn-primary">保存配置</button>
      </div>

      <div className="card">
        <h2>阿里云OSS配置</h2>

        <div className="form-group">
          <label>Endpoint</label>
          <input type="text" placeholder="oss-cn-hangzhou.aliyuncs.com" />
        </div>

        <div className="form-group">
          <label>Bucket</label>
          <input type="text" placeholder="genvideo" />
        </div>

        <div className="form-group">
          <label>Access Key</label>
          <input type="password" placeholder="Access Key" />
        </div>

        <div className="form-group">
          <label>Secret Key</label>
          <input type="password" placeholder="Secret Key" />
        </div>

        <button className="btn btn-primary">保存配置</button>
      </div>
    </div>
  )
}
