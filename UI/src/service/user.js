/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:26:47
 * @Description: User API
 * @FilePath: \CasaOS-UI\src\service\user.js
 */
import { api } from "./service.js";

const user = {
    //login
    login(data) {
        return api.post("user/login", data);
    },

    // Create UserName and Password
    createUsernameAndPaword(data) {
        return api.post("/user/setusernamepwd", data);
    },

    // Change User Avatar
    changeAvatar(data) {
        return api.post("/user/changhead", data);
    },

    // Change UserName
    changeUserName(data) {
        return api.put("/user/changusername", data);
    },

    // Change User Password
    changePassword(data) {
        return api.put("/user/changuserpwd", data);
    },

    // Get user info
    getUserInfo() {
        return api.get("/user/info");
    },
    
    // Change User Info
    changeUserInfo(data) {
        return api.post('/user/changuserinfo', data)
    }
}
export default user;
