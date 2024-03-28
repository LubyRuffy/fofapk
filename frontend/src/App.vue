<script setup>
import {reactive} from 'vue'
import {FofaStat, StartTask, UpdateScore} from '../wailsjs/go/main/App'
import {BrowserOpenURL, EventsEmit, EventsOn} from '../wailsjs/runtime/runtime';
import { ElNotification } from 'element-plus'
import {Search} from '@element-plus/icons-vue'

const data = reactive({
  taskId: "", // 任务ID
  query1: "host=\"fofa.info\"", // 参赛者1提交的规则
  query2: "domain=\"fofa.info\"", // 参赛者2提交的规则
  score1: 0, // 参赛者1的分数
  score2: 0, // 参赛者2的分数
  progress: 0, // 进度条
  showQuery: false, // 是否显示规则明文
  showFromColor: false, // 是否显示来源的颜色区分
  running: false, // 是否正在运行
  logs: [], // 日志
  data: null,
})

// ==========函数定义===========


function error(errMsg) {
  data.running = false;
  ElNotification({
    title: 'Error',
    message: errMsg,
    type: 'error',
  })
}

function success(msg) {
  data.running = false;
  ElNotification({
    title: 'Success',
    message: msg,
    type: 'success',
  })
}

function openFofa(q) {
  BrowserOpenURL(`https://fofa.info/result?qbase64=${btoa('host=="'+q+'"')}`)
}

function openFofaIP(ip) {
  BrowserOpenURL(`https://fofa.info/hosts/${ip}`)
}

function fofaStatsOfIP(ip) {
  this.loading = true;
  try {
    FofaStat(ip).then(function (response){
      if (response.data != null) {
        let text = "";
        for (let i = 0; i < response.data.length; i++) {
          text += "<br/><h3>" + response.data[i].Name + "</h3><br/>";
          for (let j = 0; j < response.data[i].Items.length; j++) {
            text += response.data[i].Items[j].Name + ":"
                + response.data[i].Items[j].Count + "<br/>";
          }

        }
        document.getElementById(ip).innerHTML = text;
      }
    })
  } finally {
    this.loading = false;
  }
}

// ==========消息处理===========


EventsOn('onProgress', (response) => {
  data.progress = response.progress;
  if (response.finished) {
    success('Finished')
  }
  data.logs.push(...response.logs);
});


EventsOn('onData', (response) => {
  data.logs.push(...response.logs);
  data.data = response.data;
});

EventsOn('onError', (response) => {
  if (response.logs != null) {
    data.logs.push(...response.logs);
  }
  error(response.error);
});


// ==========事件处理===========

function load() {
  data.running = true;
  data.logs = [];
  data.data = [];
  data.score1 = 0;
  data.score2 = 0;
  StartTask(data.query1, data.query2).then(result => {
    data.taskId = result.data.taskid
    data.resultText = result.error
  })
}

const updateScore = (ip, score) => {
  data.running = true;
  data.data = [];
  try{
    UpdateScore(data.taskId, ip, score).then(result => {
      if (result.error != null && result.error.length>0) {
        error(result.error)
      } else {
        success('update score successfully')
        data.score1 = result.data.score1
        data.score2 = result.data.score2
        data.data = result.data.data
      }
    })
  } finally {
    data.running = false;
  }
}


// ==========样式处理===========
const tableRowClassName = (a) => {
  if (!data.showFromColor) {
    return ''
  }
  if (a.row.from === 1) {
    return 'warning-row'
  } else if (a.row.from === 2) {
    return 'success-row'
  }
  return ''
}


</script>

<template>
  <el-container style="height: 100vh;">
    <el-header><h1>FOFA PK台</h1></el-header>
    <el-main>
      <span>{{data.taskId}}</span>
      <el-progress :percentage="data.progress" v-if="data.running"/>
      <div id="input" class="input-box">
        <el-row>
          <el-col :span="10" :style="{boxShadow: 'var(--el-box-shadow-lighter)'}">
            FOFA工程师 A <br/>
            <el-input v-model="data.query1"
                      :type="data.showQuery ? 'text' : 'password'"
                      style="width: 100%" :disabled="data.running"/>
            <br/> 得分: <el-tag
              :type="data.score1>=data.score2?'primary':'danger'"
              :size="data.score1>=data.score2?'large':'small'"
          >{{ data.score1 }}</el-tag>
          </el-col>
          <el-col :span="4">
            <el-button @click="load" :disabled="data.running">Load</el-button>
            <br/>显示明文<el-switch v-model="data.showQuery" />
            <br/>颜色区分<el-switch v-model="data.showFromColor" />
<!--            <br/>显示相同数据<el-switch v-model="data.showEqual" />-->
          </el-col>
          <el-col :span="10" :style="{boxShadow: 'var(--el-box-shadow-lighter)'}">
            FOFA工程师 B <br/>
            <el-input v-model="data.query2"
                      :type="data.showQuery ? 'text' : 'password'"
                      style="width: 100%" :disabled="data.running"/>
            <br/> 得分: <el-tag
              :type="data.score2>=data.score1?'primary':'danger'"
              :size="data.score2>=data.score1?'large':'small'"
          >{{ data.score2 }}</el-tag>
          </el-col>
        </el-row>
        <el-row>
          <el-skeleton :rows="5" animated v-if="data.running"/>
          <el-table :data="data.data" border style="width: 100%" :row-class-name="tableRowClassName" v-if="!data.running">
            <el-table-column type="expand">
              <template #default="props">
                <el-row style="margin: 1rem;">
                  <el-col :span="12">
                    <p m="t-0 b-2">Host: <a href="#" @click="openFofa(props.row.host)">{{ props.row.host }}</a></p>
                    <p m="t-0 b-2">IP: <a href="#" @click="openFofaIP(props.row.ip)">{{ props.row.ip }}</a></p>
                    <p m="t-0 b-2">Port: {{ props.row.port }}</p>
                    <p m="t-0 b-2">Domain: {{ props.row.domain }}</p>
                    <p m="t-0 b-2">Cert: {{ props.row.certs_subject_cn }}</p>
                    <p m="t-0 b-2">Title: {{ props.row.title }}</p>
                  </el-col>
                  <el-col :span="12">
                    <el-button :icon="Search" @click="fofaStatsOfIP(props.row.ip)">获取统计数据</el-button>
                    <div :id="props.row.ip"></div>
                  </el-col>
                </el-row>
              </template>
            </el-table-column>
            <el-table-column prop="ip" label="IP" width="180" />
            <el-table-column prop="title" label="Title" />
            <el-table-column label="Operations" width="250">
              <template #default="scope">
                <el-button size="small" type="success" @click="updateScore(scope.row.ip, 1)">
                  Valid
                </el-button>
                <el-button size="small" type="danger" @click="updateScore(scope.row.ip, -1)">
                  Invalid
                </el-button>
                <el-button size="small" @click="updateScore(scope.row.ip, 0)">
                  Unknown
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-row>
      </div>
    </el-main>
    <el-footer class="scrollable-footer" style="margin-top: auto; height: 100px; overflow-y: auto; border: 1px solid #ccc; text-align: left;">
      <div v-for="(log, index) in data.logs" :key="index" class="footer-log-item">
        {{ log }}
      </div>
    </el-footer>
  </el-container>
</template>

<style>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}
.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-9);
}

.scrollable-footer {
  /* 添加内边距，让内容与边框保持一定距离 */
  padding: 10px;

  /* 确保滚动条在内容区域内部而不是溢出边界 */
  box-sizing: border-box;

  /* 设置滚动条样式（仅适用于Webkit内核浏览器，可补充其他内核的滚动条样式） */
  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.3);
    border-radius: 4px;
  }
}
.footer-log-item {
  /* 避免每条日志之间的空白 */
  margin-bottom: 0;
}
</style>
