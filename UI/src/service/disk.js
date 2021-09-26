/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:26:02
 * @Description: Disk API
 * @FilePath: \CasaOS-UI\src\service\disk.js
 */
import { api } from "./service.js";

const disk = {
    // get Path list
    diskInfo() {
        return api.get('/disk/info');
    },
    diskList() {
        return api.get('/disk/list');
    },
    // System path
    renamePath(oldpath, path) {
        let data = {
            oldpath: oldpath,
            newpath: path
        }
        return api.get('/zima/rename', data);
    },
    // Make a new Dir
    mkdir(path) {
        let data = {
            path: path
        }
        return api.get('/zima/mkdir', data)
    }
}
export default disk;
