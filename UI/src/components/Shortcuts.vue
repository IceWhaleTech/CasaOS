<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:23:49
 * @Description: 
 * @FilePath: \CasaOS-UI\src\components\Shortcuts.vue
-->
<template>
  <div class="has-text-left mt-6">
    <div class="title-bar is-flex is-align-items-center">
      <h1 class="title is-4  has-text-white is-flex-shrink-1">Shortcuts</h1>
      <div class="buttons ">
        <b-button icon-left="plus" type="is-dark" size="is-small" rounded>Add shortcut</b-button>
      </div>
    </div>

    <div class="columns is-variable is-2 is-multiline ">
      <div class="column is-narrow is-3" v-for="n in 10" :key="n">
        <div class="wuji-card is-flex is-align-items-center ">

          <figure class="image is-32x32 simg">
            <img :src="require('@/assets/img/icon.png')">
          </figure>
          <p class="ml-4 flex1 one-line">Test</p>
          <div class="action-btn1">
            <b-dropdown aria-role="list" position="is-bottom-left" append-to-body>
              <template #trigger>
                <p role="button">
                  <b-icon pack="fas" icon="ellipsis-v" size="is-small">
                  </b-icon>
                </p>
              </template>

              <b-dropdown-item aria-role="listitem">Action</b-dropdown-item>
              <b-dropdown-item aria-role="listitem">Another action</b-dropdown-item>
              <b-dropdown-item aria-role="listitem">Something else</b-dropdown-item>
            </b-dropdown>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  name: 'shortcuts'

}
</script>

<style>
</style>