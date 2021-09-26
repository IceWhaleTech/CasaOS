<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-23 18:21:13
 * @Description: Top bar 
 * @FilePath: /CasaOS-UI/src/components/TopBar.vue
-->

<template>
  <div class="navbar top-bar is-flex is-align-items-center">
    <div class="navbar-brand ml-3">
      <b-dropdown aria-role="list" class="navbar-item" @active-change="onOpen">
        <template #trigger>
          <p role="button">
            <b-icon pack="fas" icon="sliders-h">
            </b-icon>
          </p>
        </template>

        <b-dropdown-item aria-role="menu-item" :focusable="false" custom>
          <h2 class="title is-4">CasaOS Setting</h2>
          <div class="is-flex is-align-items-center item">
            <div class="is-flex is-align-items-center flex1">
              <b-icon pack="fas" icon="sync-alt" class="mr-1"></b-icon> <b>Update</b>
            </div>
            <div>
              v{{updateInfo.current_version}}
            </div>
            <!-- <b-field>
              <b-switch v-model="barData.auto_update" type="is-dark" size="is-small" class="is-flex-direction-row-reverse mr-0" @input="saveData">
                Auto-Check
              </b-switch>
            </b-field> -->
          </div>

          <div class="is-flex is-align-items-center pl-5" v-if="!updateInfo.is_need">
            {{latestText}}
            <b-icon type="is-success" pack="fas" icon="check" class="ml-1"></b-icon>
          </div>
          <div class="is-flex is-align-items-center is-justify-content-end update-container pl-5" v-if="updateInfo.is_need">
            <div class="flex1">{{updateText}}</div>
            <b-button type="is-dark" size="is-small" class="ml-2" :loading="isUpdating" rounded @click="updateSystem">Update</b-button>
          </div>

        </b-dropdown-item>
      </b-dropdown>

    </div>

    <div class="navbar-menu">
      <div class="navbar-end mr-3">
        <!-- <b-icon pack="far" icon="comment-alt"></b-icon> -->
      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: "top-bar",
  data() {
    return {
      timer: 0,
      barData: {
        auto_update: false,
        background: "",
        background_type: "d",
        search_engine: "google",
        search_switch: false,
        shortcuts_switch: false,
        widgets_switch: false
      },
      updateInfo: {
        current_version: '0',
        is_need: false,
        version: Object
      },
      isUpdating: false,
      latestText: "Currently the latest version",
      updateText: "A new version is available!"
    }
  },
  created() {
    this.getConfig();
  },

  methods: {
    /**
     * @description: Get CasaOs Configs
     * @return {*} void
     */
    getConfig() {
      this.$api.info.systemConfig().then(res => {
        if (res.data.success == 200) {
          this.barData = res.data.data
        }
      })
    },

    /**
     * @description: Save CasaOs Configs
     * @return {*} void
     */
    saveData() {
      this.$api.info.saveSystemConfig(this.barData).then(res => {
        if (res.data.success == 200) {
          console.log(res);
        }
      })
    },

    /**
     * @description: Handle Dropmenu state
     * @param {Boolean} isOpen
     * @return {*} void
     */
    onOpen(isOpen) {
      if (isOpen) {
        this.$api.info.checkVersion().then(res => {
          if (res.data.success == 200) {
            this.updateInfo = res.data.data
          }
        })
      }
    },

    /**
     * @description: Update System Version and check update state
     * @return {*} void
     */
    updateSystem() {
      this.isUpdating = true;
      this.$api.info.updateSystem().then(res => {
        if (res.data.success == 200) {
          console.log(res.data.data);
        }
      });
      this.checkUpdateState();
    },
    /**
     * @description: check update state if is_need is false then reload page
     * @return {*} void
     */
    checkUpdateState() {
      this.timer = setInterval(() => {
        this.$api.info.checkVersion().then(res => {
          if (res.data.success == 200) {
            if (!res.data.data.is_need) {
              clearInterval(this.timer);
              location.reload();
            }
          }
        })
      }, 3000)
    }
  },

}
</script>
