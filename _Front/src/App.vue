<template>
    <div style="display: flex;justify-content: center;gap: 30px;">
        <div style="width: 55%;">
            <div style="margin: 16px 0;display: flex;gap: 20px;">
                <el-input v-model="currentPath" disabled></el-input>
                <el-upload
                    action="fileMgr/upload"
                    :data="{ path: currentPath }"
                    :show-file-list="false"
                    :before-upload="beforeUpload"
                    :on-success="handleSuccess"
                    :on-error="handleError"
                >
                    <el-button type="primary" :loading="uploadLoading">上传</el-button>
                </el-upload>
            </div>
            <el-table :data="tableData" @row-dblclick="routePath" height="600">
                <el-table-column prop="isDir" width="50">
                    <template #default="{ row }">
                        <div style="height: 100%;display: flex;flex-direction: column;justify-content: center;align-items: center;">
                            <i-ep-folder v-if="row.isDir" />
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="Name" />
                <el-table-column prop="size" label="Size" width="80">
                    <template #default="{ row }"> 
                        <span>{{ row.size != 0 ? prettyBytes(row.size) : "" }}</span>
                    </template>
                </el-table-column>
                <el-table-column label="Action" width="80" align="center">
                    <template #default="{ row }">
                        <el-button type="danger" circle size="small" :disabled="row.isDir" @click="routePath(row, 'del')">
                            <i-ep-delete #icon />
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div style="width: 35%;margin: 16px 0;">
            <el-input v-model="sendText" type="textarea" :rows="31" placeholder="请输入内容" />
            <div>
                <el-button type="primary" size="small" :loading="sendLoading" @click="handleSend">发送到剪贴板</el-button>
                <el-button type="primary" size="small" :loading="sendLoading" @click="handleReceive">获取自剪贴板</el-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import axios from 'axios'
import prettyBytes from 'pretty-bytes'
import { onMounted, reactive, ref } from 'vue'

let tableData = reactive([])
onMounted(_ => getList("D:/"))
let currentPath = ref("")
let currentPathArr = ["D:"]
let routePath = (row, del) => {
    if (row.name == "..") {
        currentPathArr.pop()
    } else {
        currentPathArr.push(row.name)
    }
    let path = currentPathArr.length == 1 ? currentPathArr[0] + "/" : currentPathArr.join("/")
    if (del === "del") {
        currentPathArr.pop()
        handelDelete(path)
        return
    }
    if (row.isDir) {
        getList(path)
    } else {
        currentPathArr.pop()
        getFile(path, row)
    }
}
let getFile = (path, row) => {
    if (row.size > 1000 * 1000 * 30) {
        ElMessage({
            message: 'Size > ' + prettyBytes(1000 * 1000 * 30),
            type: 'warning',
        })
    } else {
        const link = document.createElement('a')
        link.style.display = "none"
        link.href = "fileMgr/getFile?fileName=" + path
        link.setAttribute('download', row.name)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    }
}
let handelDelete = (path) => {
    const options = {
        method: 'POST',
        url: 'fileMgr/delFile',
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
        data: { path: path }
    }
    axios.request(options).then(res => {
        if (res.status == 200 && res.data.succMsg == 'ok') {
            getList(currentPath.value)
        }
    })
}
let getList = (path) => {
    tableData.length = 0
    const options = {
        method: 'POST',
        url: 'fileMgr/getList',
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
        data: { path: path }
    }
    axios.request(options).then(res => {
        if (res.status == 200 && res.data.succMsg == 'ok') {
            currentPath.value = path
            if (path != "D:/") {
                tableData.push({ name: "..", isDir: true, size: '' })
            }
            tableData.push(...res.data.data)
        }
    })
}

let uploadLoading = ref(false)
let beforeUpload = (file) => {
    if (file.size > 1000 * 1000 * 30) {
        ElMessage({
            message: 'Size > ' + prettyBytes(1000 * 1000 * 30),
            type: 'warning',
        })
        return false
    }
    uploadLoading.value = true
    return true
}
let handleSuccess = (response) => {
    if (response.succMsg == 'ok') {
        uploadLoading.value = false
        getList(currentPath.value)
    }
}
let handleError = (err) => {
    uploadLoading.value = false
    ElMessage({
        message: err,
        type: 'warning',
    })
}

let sendText = ref('')
let sendLoading = ref(false)
let handleSend = () => {
    sendLoading.value = true
    const options = {
        method: 'POST',
        url: 'tools/SendText',
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
        data: { text: sendText.value }
    }
    axios.request(options).then(res => {
        if (res.status == 200 && res.data.succMsg == 'ok') {
            ElMessage({
                message: "发送成功",
                type: 'success',
            })
            sendText.value = ''
        }
        sendLoading.value = false
    }).catch(err =>{
        sendLoading.value = false
    })
}
let handleReceive = () => {
    sendLoading.value = true
    const options = {
        method: 'POST',
        url: 'tools/GetText',
        headers: { 'content-type': 'application/x-www-form-urlencoded' }
    }
    axios.request(options).then(res => {
        if (res.status == 200 && res.data.succMsg == 'ok') {
            sendText.value = res.data.data
        }
        sendLoading.value = false
    }).catch(err =>{
        sendLoading.value = false
    })
}
</script>

<style scoped></style>
