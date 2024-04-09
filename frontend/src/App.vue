<script setup>
import {reactive, ref} from 'vue'
import {FofaStat, StartTask, UpdateConfig, UpdateScore} from '../wailsjs/go/main/App'
import {BrowserOpenURL, EventsEmit, EventsOn} from '../wailsjs/runtime/runtime';
import { ElNotification, ElTree  } from 'element-plus'
import {Flag, Search} from '@element-plus/icons-vue'
import FofaHeader from "./components/FofaHeader.vue";

const data = reactive({
  taskId: "", // 任务ID
  fofaEngineerA: {
    name: "FOFA工程师 A", // 参赛者的名字
    query: "host=\"fofa.info\"", // 参赛者提交的规则
    score: 0, // 参赛者的分数
    size: 0, // 参赛者的总数
  },
  fofaEngineerB: {
    name: "FOFA工程师 B", // 参赛者的名字
    query: "domain=\"fofa.info\"", // 参赛者提交的规则
    score: 0, // 参赛者的分数
    size: 0, // 参赛者的总数
  },
  diffSize: 0,
  progress: 0, // 进度条
  showFromColor: false, // 是否显示来源的颜色区分
  running: false, // 是否正在运行
  logs: [], // 日志
  data: null, // 原始数据
  treeData: null, // 聚类数据
  statsLoading: false,
  activePage: "all",
  fofaKey: "",
  selected: false, // group 是否有选中
  fetchSize: 2000, // 每次获取的条数
})

const treeRef = ref();

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

function fofaStatsOfIP(row) {
  let ip = row.ip;
  data.statsLoading = true;

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
    // data.statsLoading = false;
  }).finally(()=>{
    data.statsLoading = false;
  })
}

function mapToArraySortBySize(mapData) {
  return Array.from(mapData.values()).sort((a, b) => b.children.length - a.children.length);
}

function addEntryToMap(map, key, value) {
  if (!map.has(key)) {
    map.set(key, {
      label: key,
      children: [],
    });
  }
  map.get(key).children.push(value);
}


function treeDataFromRecords(data) {
  const domains = new Map()
  const titles = new Map()
  const as_organizations = new Map()
  const fids = new Map()
  data.forEach(item => {
    let rowItem = {
      label: item.host,
      raw: item
    }
    if (item.domain!=null) {
      addEntryToMap(domains, item.domain, rowItem);
    }

    if (item.title!=null) {
      addEntryToMap(titles, item.title, rowItem);
    }

    if (item.as_organization!=null) {
      addEntryToMap(as_organizations, item.as_organization, rowItem);
    }

    if (item.fid!=null) {
      addEntryToMap(fids, item.fid, rowItem);
    }
  });
  return [
    {
      label: 'Domain',
      children: mapToArraySortBySize(domains),
    },
    {
      label: 'Title',
      children: mapToArraySortBySize(titles),
    },
    {
      label: 'ASOrg',
      children: mapToArraySortBySize(as_organizations),
    },
    {
      label: 'FID',
      children: mapToArraySortBySize(fids),
    },
  ]
}

const onDataUpdate = (respData) => {
  data.data = respData.map(item => {
    return {
      loading: false,
      ...item
    }
  });
  data.treeData = treeDataFromRecords(data.data)
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

  onDataUpdate(response.data);

  data.fofaEngineerA.size = response.size1;
  data.fofaEngineerB.size = response.size2;
  data.diffSize = response.diffSize;
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
  data.fofaEngineerA.score = 0;
  data.fofaEngineerA.size = 0;
  data.fofaEngineerB.score = 0;
  data.fofaEngineerB.size = 0;
  data.diffSize = 0;
  StartTask(data.fofaEngineerA.query, data.fofaEngineerB.query).then(result => {
    data.taskId = result.data.taskid
    data.resultText = result.error
  })
}

const updateScore = (ips, score) => {
  data.running = true;
  data.data = [];

  UpdateScore(data.taskId, ips, score).then(result => {
    if (result.error != null && result.error.length>0) {
      error(result.error)
    } else {
      success('update score successfully')
      data.fofaEngineerA.score = result.data.score1
      data.fofaEngineerB.score = result.data.score2
      data.diffSize = result.data.diffSize
      onDataUpdate(result.data.data);
      if (result.data.logs != null) {
        data.logs.push(...result.data.logs);
      }
    }

  }).finally(()=>{
    data.running = false;
  })
}

const updateScoreBySelected = (score) => {
  let ips = []
  treeRef.value.getCheckedNodes(true).forEach((item) => {
    if (item.raw != null) {
      ips.push(item.raw.ip)
    }
  })
  updateScore(ips, score)
}

const updateConfig = () => {
  UpdateConfig(data.fofaKey, parseInt(data.fetchSize, 10)).then(result => {
    if (result.error != null && result.error.length>0) {
      error(result.error)
    } else {
      success('update config successfully')
    }
  }).finally(()=>{
  })
}

const handleCheck = (
    tree,
    checked,
    indeterminate
) => {
  data.selected = (checked['checkedNodes'].length > 0)
  console.log(data, checked, indeterminate)
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
    <el-header><h1>FOFA PK台 <span style="font-size: 12px">v0.2</span></h1></el-header>
    <el-main>
      <FofaHeader :data="data" :load="load"/>
      <el-row>
        <el-skeleton :rows="5" animated v-if="data.running"/>
        <el-tabs model-value="all" v-if="!data.running" style="width: 100%" type="card">
          <el-tab-pane name="all">
            <template #label>
              <el-badge :value="data.diffSize" :max="999" class="item" type="primary">
                差异数据列表
              </el-badge>
            </template>
            <el-table :data="data.data" border style="width: 100%" :row-class-name="tableRowClassName" v-if="!data.running">
              <el-table-column type="expand">
                <template #default="props">
                  <el-row style="margin: 1rem;">
                    <el-col :span="12">
                      <p m="t-0 b-2">Host: <a href="#" @click="openFofa(props.row.host)">{{ props.row.host }}</a></p>
                      <p m="t-0 b-2">IP: <a href="#" @click="openFofaIP(props.row.ip)">{{ props.row.ip }}</a></p>
                      <p m="t-0 b-2">Port: {{ props.row.port }}</p>
                      <p m="t-0 b-2">Protocol: {{ props.row.protocol }}</p>
                      <p m="t-0 b-2">Domain: {{ props.row.domain }}</p>
                      <p m="t-0 b-2">Cert: {{ props.row.certs_subject_cn }}</p>
                      <p m="t-0 b-2">Title: {{ props.row.title }}</p>
                      <p m="t-0 b-2">FID: {{ props.row.fid }}</p>
                    </el-col>
                    <el-col :span="12">
                      <el-button :icon="Search" @click="fofaStatsOfIP(props.row)" :loading="data.statsLoading">获取统计数据</el-button>
                      <div :id="props.row.ip"></div>
                    </el-col>
                  </el-row>
                </template>
              </el-table-column>
              <el-table-column prop="ip" label="IP" width="180" />
              <el-table-column prop="title" label="Title" />
              <el-table-column label="Operations" width="250">
                <template #default="scope">
                  <el-button size="small" type="success" @click="updateScore([scope.row.ip], 1)">
                    Valid
                  </el-button>
                  <el-button size="small" type="danger" @click="updateScore([scope.row.ip], -1)">
                    Invalid
                  </el-button>
                  <el-button size="small" @click="updateScore([scope.row.ip], 0)">
                    Unknown
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="分组" name="group">
            <el-button size="small" type="success" @click="updateScoreBySelected(1)" :disabled="!data.selected">
              Valid
            </el-button>
            <el-button size="small" type="danger" @click="updateScoreBySelected(-1)" :disabled="!data.selected">
              Invalid
            </el-button>
            <el-button size="small" @click="updateScoreBySelected(0)" :disabled="!data.selected">
              Unknown
            </el-button>
            <el-tree
                ref="treeRef"
                style="width: 100%"
                :data="data.treeData"
                show-checkbox
                node-key="label"
                :default-expanded-keys='["Domain","Title","ASOrg","FID"]'
                @check="handleCheck"
            />
          </el-tab-pane>
          <el-tab-pane label="配置项" name="config">
            FOFA_KEY（如何环境变量中进行了配置，这里可以不填，默认会调用系统的环境变量配置）: <el-input placeholder="FOFA_KEY" v-model="data.fofaKey"></el-input>
            FOFA请求获取的数据量大小: <el-input placeholder="fofa size" v-model="data.fetchSize"></el-input>
            <el-button @click="updateConfig()">Update</el-button>
          </el-tab-pane>
        </el-tabs>


      </el-row>
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
