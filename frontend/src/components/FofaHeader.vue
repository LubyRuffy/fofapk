<!--
  - Copyright (c) 2024. LubyRuffy. All rights reserved.
  -->

<script setup>
import FofaEngineer from "./FofaEngineer.vue";
import {defineProps, reactive} from 'vue';

// 定义组件的props
const props = defineProps({
  data: Object,
  load: Function,
});

const config = reactive({
  showQuery: false, // 是否显示规则明文
})

</script>

<template>
  <!--      <span>{{data.taskId}}</span>-->
  <el-progress :percentage="data.progress" v-if="data.running"/>
  <el-row>
    <el-col :span="10">
      <FofaEngineer :engineer="data.fofaEngineerA"
                    :running="data.running"
                    :show-query="config.showQuery"
                    :winner="data.fofaEngineerA.score>data.fofaEngineerB.score"
      />
    </el-col>
    <el-col :span="4">
      <el-button @click="load" :disabled="data.running">Load</el-button>
      <br/>显示明文<el-switch v-model="config.showQuery" />
      <br/>颜色区分<el-switch v-model="data.showFromColor" />
      <!--            <br/>显示相同数据<el-switch v-model="data.showEqual" />-->
    </el-col>
    <el-col :span="10">
      <FofaEngineer :engineer="data.fofaEngineerB"
                    :running="data.running"
                    :show-query="config.showQuery"
                    :winner="data.fofaEngineerB.score>data.fofaEngineerA.score"
      />
    </el-col>
  </el-row>
</template>

<style scoped>

</style>