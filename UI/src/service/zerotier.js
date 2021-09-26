/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:26:50
 * @Description: Zerotier API
 * @FilePath: \CasaOS-UI\src\service\zerotier.js
 */
import { api } from "./service.js";

const zerotier = {
    //Check if Need login to zerotier
    isLogin() {
        return api.get("/zerotier/islogin");
    },
    //Login
    login(data) {
        return api.post("/zerotier/login", data);
    },
    //Register
    register(data) {
        return api.post('/zerotier/register', data);
    },
    //networklist
    networkLits() {
        return api.get('/zerotier/list'); 
    },
    //joinNetwork
    joinNetwork(id) {
        return api.post(`/zerotier/join/${id}`);
    },
    // leaveNetwork
    leaveNetwork(id) {
        return api.post(`/zerotier/leave/${id}`);
    },
    // Get Network detial
    networkDetail(id) {
        return api.get(`/zerotier/info/${id}`);
    },
    // Edit Network
    editNetwork(id, data) {
        return api.put(`/zerotier/edit/${id}`, data)
    },
    // Delete A Network
    delNetwork(id) {
        return api.delete(`/zerotier/network/${id}/del`)
    },
    createNetwork() {
        return api.post('/zerotier/create')
    },
    // Get Network member list
    getMembers(id) {
        return api.get(`/zerotier/member/${id}`)
    },
    // Edit Member
    editMember(id, mId, data) {
        return api.put(`/zerotier/member/${id}/edit/${mId}`, data)
    },
    // Delete Member
    delMemeber(id, mId) {
        return api.delete(`/zerotier/member/${id}/del/${mId}`)
    }
}
export default zerotier;
