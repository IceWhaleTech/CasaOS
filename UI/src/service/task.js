/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:26:45
 * @Description: Task API
 * @FilePath: \CasaOS-UI\src\service\task.js
 */
import { api } from "./service.js";

const task = {
    //List
    list() {
        return api.get("/task/list");
    },
    //Mark
    completion(id) {
        return api.put(`/task/completion/${id}`);
    }
}
export default task; 
