<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-18 23:13:58
 * @Description: The left bottom brand bar
 * @FilePath: \CasaOS-UI\src\components\BrandBar.vue
-->

<template>
  <div class="brand-bar is-flex is-align-items-center has-text-white">
    <figure class="image is-32x32">
      <img :src="require('@/assets/img/casa.svg')">
    </figure>
    <span class="is-size-4 mr-3 ml-3">CasaOS</span>
    <span>Made by IceWhale with ❤ and you !️</span>
  </div>
</template>

<script>
export default {
  name: "brand-bar"
}
</script>

<style>
</style>