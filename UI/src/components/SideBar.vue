<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-18 23:20:56
 * @Description: 
 * @FilePath: \CasaOS-UI\src\components\SideBar.vue
-->
<template>
  <div class="side-bar mr-5">
    <clock></clock>
  </div>
</template>

<script>
import Clock from '../assets/widgets/Clock.vue'


export default {
  name: 'side-bar',
  components: { Clock },

}
</script>

<style>
</style>