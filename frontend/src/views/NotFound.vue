<template>
  <div class="not-found">
    <div class="content">
      <div class="error-code">404</div>
      <h1>页面未找到</h1>
      <p class="description">抱歉，您访问的页面不存在或已被移除。</p>
      <div class="actions">
        <router-link to="/" class="home-btn">返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
}

.content {
  text-align: center;
  max-width: 500px;
}

.error-code {
  font-size: 8rem;
  font-weight: bold;
  color: #00ff00;
  line-height: 1;
  margin-bottom: 1rem;
  text-shadow: 0 0 30px rgba(0, 255, 0, 0.3);
}

h1 {
  color: #fff;
  font-size: 2rem;
  margin-bottom: 1rem;
}

.description {
  color: #888;
  font-size: 1.1rem;
  margin-bottom: 2rem;
}

.actions {
  display: flex;
  justify-content: center;
}

.home-btn {
  background: #00ff00;
  color: #000;
  padding: 0.75rem 2rem;
  border-radius: 6px;
  text-decoration: none;
  font-weight: bold;
  transition: all 0.3s;
}

.home-btn:hover {
  background: #00cc00;
  transform: translateY(-2px);
}
</style>