<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-18 23:14:33
 * @Description: The right bottom contact bar
 * @FilePath: \CasaOS-UI\src\components\ContactBar.vue
-->

<template>
  <div class="contact-bar is-flex is-align-items-center has-text-white pl-3 pr-3">
    <a href="#">
      <b-icon pack="fas" icon="map-signs" size=""></b-icon>
    </a>
    <a href="https://discord.gg/Gx4BCEtHjx" target="_blank">
      <b-icon pack="fab" icon="discord" size=""></b-icon>
    </a>
    <a href="https://github.com/ZimaBoard/CasaOS" target="_blank">
      <b-icon pack="fab" icon="github" size=""></b-icon>
    </a>
  </div>
</template>

<script>
export default {
  name: "contact-bar"
}
</script>
