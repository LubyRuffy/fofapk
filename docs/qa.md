# 问题记录

## vue

### vue中reactive和ref有什么不同？
在 Vue 3 中，reactive 和 ref 都是用来创建响应式数据的不同方法，它们的主要区别在于处理数据结构的方式以及如何访问和更新数据：
- reactive
  - reactive 函数用于创建一个响应式的**对象**，它可以深度监控对象的所有属性变化。
  - 当调用 reactive 时，它会返回一个新的代理对象，该对象与原始对象具有相同的属性和结构，但所有属性都成为响应式的。
  - 在模板或计算属性中可以直接访问和修改 reactive 对象的属性，无需在其后附加 .value。
- ref
  - ref 函数主要用于创建一个响应式的**单个值**包装器，它返回一个带有 .value 属性的对象，该属性就是你要跟踪的值。
  - 要访问或修改 ref 包装的值，你需要通过 .value 属性进行操作。
我们通常用一个data的object来存储数据，直接用reactive就好了。

### 如何绑定一个值到button的enabled状态
```vue
<button disabled="!data.running">Load</button>
```

### 如何绑定一个值到input
```vue
<input type="text" v-model="message">
```

### 如何通过div展示一个数组
```vue
<div v-for="(log, index) in data.logs" :key="index">
  {{ log }}
</div>
```

## element-ui

### input替换
```vue
<el-input v-model="data.query1"></el-input>
```

### input动态修改password类型
```vue
<el-input v-model="data.query1" :type="data.showQuery ? 'text' : 'password'"></el-input>
```

### table做成可以折叠展开的方式
用`expand`：
```vue
<el-table :data="data.data" border style="width: 100%">
  <el-table-column type="expand">
    <template #default="props">
      <div m="4">
        <p m="t-0 b-2">Host: {{ props.row.host }}</p>
        <p m="t-0 b-2">Domain: {{ props.row.domain }}</p>
        <p m="t-0 b-2">Cert: {{ props.row.certs_subject_cn }}</p>
      </div>
    </template>
  </el-table-column>
  <el-table-column prop="ip" label="IP" width="180" />
  <el-table-column prop="title" label="Title" />
</el-table>
```

## javascript
### javascript把一个数组拼接到另一个数组中去
可以先使用 Spread Operator(...) 将数组展开，然后再使用 push() 方
```javascript
data.logs.push(...response.logs);
```
