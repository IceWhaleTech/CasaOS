/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-23 17:41:01
 * @Description: System HardWare Info API
 * @FilePath: /CasaOS-UI/src/service/info.js
 */
import { api } from "./service.js";

const info = {
    //CPU info
    cpuInfo() {
        return api.get("/zima/getcpuinfo");
    },
    //Memory Info
    memoryInfo() {
        return api.get("/zima/getmeminfo");
    },
    //Network Info
    networkInfo() {
        return api.get('/zima/getnetinfo');
    },
    //Disk Info
    diskInfo() {
        return api.get('/zima/getdiskinfo');
    },
    //All Info
    allInfo() {
        return api.get('/zima/getinfo');
    },
    // System Info
    systemInfo() {
        return api.get('/zima/sysinfo');
    },
    //Get CasaOS Config
    systemConfig() {
        return api.get('/sys/config')
    },
    //Save CasaOs Config
    saveSystemConfig(data) {
        return api.post('/sys/config', data)
    },
    // Check Verison
    checkVersion() {
        return api.get('/sys/check');
    },
    //Update System
    updateSystem(){
        return api.post('/sys/update');
    }
}
export default info;
