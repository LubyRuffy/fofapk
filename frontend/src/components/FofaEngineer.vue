<!--
  - Copyright (c) 2024. LubyRuffy. All rights reserved.
  -->

<script setup>
import {Flag} from "@element-plus/icons-vue";
import { defineProps } from 'vue';

// 定义组件的props
const props = defineProps({
  engineer: Object,
  showQuery: Boolean,
  winner: Boolean,
  running: Boolean,
});

</script>

<template>
  <el-row>
    <el-col>
      {{engineer.name}} <br/>
      <el-input v-model="engineer.query"
                :type="showQuery ? 'text' : 'password'"
                style="width: 100%" :disabled="running"/>
    </el-col>
  </el-row>

  <el-row>
    <el-col :span="12">
      <el-statistic title="总数" :value="engineer.size" />
    </el-col>
    <el-col :span="12">
      <el-statistic :value="engineer.score" :value-style="{color: winner?'red':'green'}">
        <template #title>
          <div style="display: inline-flex; align-items: center">
            得分
            <el-icon style="margin-left: 4px" :size="12">
              <Flag />
            </el-icon>
          </div>
        </template>
      </el-statistic>
    </el-col>
  </el-row>
</template>

<style scoped>

</style>