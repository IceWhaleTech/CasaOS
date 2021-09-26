<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-22 16:27:00
 * @Description: Install Panel of Docker
 * @FilePath: /CasaOS-UI/src/components/Panel.vue
-->

<template>
  <div class="modal-card">
    <!-- Modal-Card Header Start -->
    <header class="modal-card-head">
      <div class="flex1">
        <h3 class="title is-4 has-text-weight-normal">Create a new App manually</h3>
      </div>
      <b-button icon-left="file-import" type="is-dark" size="is-small" rounded @click="showImportPanel" v-if="currentSlide == 1 && state == 'install'">Import</b-button>
    </header>
    <!-- Modal-Card Header End -->
    <!-- Modal-Card Body Start -->
    <section class="modal-card-body">

      <section v-show="currentSlide == 1">
        <ValidationObserver ref="ob1">
          <ValidationProvider rules="required" name="Image" v-slot="{ errors, valid }">
            <b-field label="Docker Image *" :type="{ 'is-danger': errors[0], 'is-success': valid }" :message="errors">
              <b-input v-model="initData.image" placeholder="e.g.,hello-world:latest" :readonly="state == 'update'"></b-input>
              <!-- <b-autocomplete :data="data" placeholder="e.g. hello-world:latest" field="image" :loading="isFetching" @typing="getAsyncData" @select="option => selected = option" v-model="initData.image" :readonly="state == 'update'"></b-autocomplete> -->
            </b-field>
          </ValidationProvider>
          <ValidationProvider rules="required" name="Name" v-slot="{ errors, valid }">
            <b-field label="App name *" :type="{ 'is-danger': errors[0], 'is-success': valid }" :message="errors">
              <b-input value="" v-model="initData.label" placeholder="Your custom App Name"></b-input>
            </b-field>
          </ValidationProvider>

          <b-field label="Icon URL">
            <b-input value="" v-model="initData.icon" placeholder="Your custom icon URL"></b-input>
          </b-field>

          <b-field label="Web UI">
            <p class="control">
              <span class="button is-static">{{baseUrl}}</span>
            </p>
            <b-input v-model="webui" placeholder="8080/web/index.html" expanded></b-input>
          </b-field>

          <b-field label="Network">
            <b-select v-model="initData.network_model" placeholder="Select" expanded>
              <optgroup v-for="net in networks" :key="net.driver" :label="net.driver">
                <option v-for="(option,index) in net.networks" :value="option.id" :key="option.name+index">
                  {{ option.name}}
                </option>
              </optgroup>
            </b-select>
          </b-field>

          <ports v-model="initData.ports" :showHostPost="showHostPort" v-if="showPorts"></ports>
          <input-group v-model="initData.volumes" label="Data Volumes" message="No App Data Volumes now, Click “+” to add one."></input-group>
          <input-group v-model="initData.envs" label="Environment Variables" message="No environment variables now, Click “+” to add one." name1="Key" name2="Value"></input-group>
          <input-group v-model="initData.devices" label="Devices" message="No devices now, Click “+” to add one."></input-group>
          <b-field label="Memory Limit">
            <vue-slider :min="256" :max="totalMemory" v-model="initData.memory"></vue-slider>
          </b-field>

          <b-field label="CPU Shares">
            <b-select v-model="initData.cpu_shares" placeholder="Select" expanded>
              <option value="10">Low</option>
              <option value="50">Medium</option>
              <option value="90">High</option>
            </b-select>
          </b-field>

          <b-field label="Restart Policy">
            <b-select v-model="initData.restart" placeholder="Select" expanded>
              <option value="on-failure">on-failure</option>
              <option value="always">always</option>
              <option value="unless-stopped">unless-stopped</option>
            </b-select>
          </b-field>
          <b-field label="App Description">
            <b-input v-model="initData.description"></b-input>
          </b-field>
          <b-loading :is-full-page="false" v-model="isLoading" :can-cancel="false"></b-loading>
        </ValidationObserver>
      </section>
      <section v-show="currentSlide == 2">
        <div class="installing-warpper">
          <lottie-animation path="./ui/img/ani/rocket-launching.json" :autoPlay="true" :width="200" :height="200"></lottie-animation>
          <h3 class="title is-6 has-text-centered" :class="{'has-text-danger':errorType == 3,'has-text-black':errorType != 3}" v-html="installText"></h3>
        </div>
      </section>

    </section>
    <!-- Modal-Card Body End -->

    <!-- Modal-Card Footer Start-->
    <footer class="modal-card-foot is-flex is-align-items-center">
      <div class="flex1"></div>
      <div>
        <b-button v-if="currentSlide == 1" :label="cancelButtonText" @click="$emit('close')" rounded />
        <b-button v-if="currentSlide == 2 && errorType == 3 " label="Back" @click="prevStep" rounded />
        <b-button v-if="currentSlide == 1 && state == 'install'" label="Install" type="is-dark" @click="installApp()" rounded />
        <b-button v-if="currentSlide == 1 && state == 'update'" label="Update" type="is-dark" @click="updateApp()" rounded />
        <b-button v-if="currentSlide == 2" :label="cancelButtonText" type="is-dark" @click="$emit('close')" rounded />
      </div>
    </footer>
    <!-- Modal-Card Footer End -->
  </div>
</template>

<script>
import axios from 'axios'
import InputGroup from './forms/InputGroup.vue';
import Ports from './forms/Ports.vue'
import ImportPanel from './forms/ImportPanel.vue'
import LottieAnimation from "lottie-vuejs/src/LottieAnimation.vue";
import VueSlider from 'vue-slider-component'
import 'vue-slider-component/theme/default.css'
import { ValidationObserver, ValidationProvider } from "vee-validate";
import "@/plugins/vee-validate";
import debounce from 'lodash/debounce'

export default {
  components: {
    Ports,
    InputGroup,
    ValidationObserver,
    ValidationProvider,
    LottieAnimation,
    VueSlider
  },
  data() {
    return {
      timer: 0,
      data: [],
      isLoading: false,
      isFetching: false,
      errorType: 1,
      currentSlide: 1,
      cancelButtonText: "Cancel",
      webui: "",
      baseUrl: "",
      totalMemory: 0,
      networks: [],
      tempNetworks: [],
      networkModes: [],
      installPercent: 0,
      installText: "",

      initData: {
        port_map: "",
        cpu_shares: 10,
        memory: 300,
        restart: "always",
        label: "",
        position: true,
        index: "",
        icon: "",
        network_model: "",
        image: "",
        description: "",
        origin: "custom",
        ports: [],
        volumes: [],
        envs: [],
        devices: [],
      }

    }
  },
  props: {
    id: String,
    state: String,
    configData: Object,
    initDatas: {
      type: Object
    }
  },

  created() {
    //If it is edit, Init data
    if (this.initDatas != undefined) {
      this.initData = this.initDatas
      this.webui = this.initDatas.port_map + this.initDatas.index
    }

    //Get Max memory info form device
    this.totalMemory = Math.floor(this.configData.memory.total / 1048576);
    this.initData.memory = this.totalMemory

    //Handling network types
    this.tempNetworks = this.configData.networks;
    this.networkModes = this.unique(this.tempNetworks.map(item => {
      return item.driver
    }))
    this.networks = this.networkModes.map(item => {
      let tempitem = {}
      tempitem.driver = item
      tempitem.networks = this.tempNetworks.filter(net => {
        return net.driver == item
      })
      return tempitem
    })

    let gg = this.tempNetworks.filter(item => {
      if (item.driver == "bridge") {
        return item
      }
    })

    this.initData.network_model = gg[0].id

    // Set Front-end base url
    this.baseUrl = `${window.location.protocol}//${document.domain}:`;
  },
  computed: {

    showPorts() {
      if (this.initData.network_model.indexOf("macvlan") > -1) {
        return false
      } else {
        return true
      }
    },
    showHostPort() {
      if (this.initData.network_model.indexOf("host") > -1) {
        return false
      } else {
        return true
      }
    }
  },
  methods: {

    /**
     * @description: Process the datas before submit
     * @param {*}
     * @return {*} void
     */
    processData() {
      // GET port map and index
      if (this.webui != "") {
        let slashArr = this.webui.split("/")
        this.initData.port_map = slashArr[0]
        this.initData.index = "/" + slashArr.slice(1).join("/");
      }

      let model = this.initData.network_model.split("-");
      this.initData.network_model = model[0]
    },

    /**
     * @description: Array deduplication
     * @param {Array} arr
     * @return {Array}
     */
    unique(arr) {
      for (var i = 0; i < arr.length; i++) {
        for (var j = i + 1; j < arr.length; j++) {
          if (arr[i] == arr[j]) {
            arr.splice(j, 1);
            j--;
          }
        }
      }
      return arr;
    },

    /**
     * @description: Back to prev Step
     * @param {*}
     * @return {*} void
     */
    prevStep() {
      this.currentSlide--;
    },

    /**
     * @description: Validate form async
     * @param {Object} ref ref of component
     * @return {Boolean} 
     */
    async checkStep(ref) {
      let isValid = await ref.validate()
      return isValid
    },

    /**
     * @description: Submit datas after valid
     * @param {*}
     * @return {*} void
     */
    installApp() {
      this.checkStep(this.$refs.ob1).then(val => {
        if (val) {
          this.processData();
          this.isLoading = true;
          this.$api.app.install(this.id, this.initData).then((res) => {
            this.isLoading = false;
            if (res.data.success == 200) {
              this.currentSlide = 2;
              this.cancelButtonText = "Continue in background"
              this.checkInstallState(res.data.data)
            } else {
              //this.currentSlide = 1;
              this.$buefy.toast.open({
                message: res.data.message,
                type: 'is-warning'
              })
            }
          })
        }
      })
    },

    /**
     * @description: Check the installation process every 250 milliseconds
     * @param {String} appId
     * @return {*} void
     */
    checkInstallState(appId) {
      this.timer = setInterval(() => {
        this.updateInstallState(appId)
      }, 250)
    },

    /**
     * @description: Update the installation status to the UI
     * @param {String} appId
     * @return {*} void
     */
    updateInstallState(appId) {
      this.$api.app.state(appId).then((res) => {
        let resData = res.data.data;
        this.installPercent = resData.speed;
        this.errorType = resData.type;
        if (this.errorType == 4) {
          try {
            let info = JSON.parse(resData.message)
            let id = (info.id != undefined) ? info.id : "";
            let progress = ""
            if (info.progressDetail != undefined) {
              let progressDetail = info.progressDetail
              if (!isNaN(progressDetail.current / progressDetail.total)) {
                progress = "<br>Progress:" + String(Math.floor((progressDetail.current / progressDetail.total) * 100)) + "%"
              }
            }
            let status = info.status
            this.installText = status + ":" + id + " " + progress
          } catch (error) {
            console.log(error);
          }
        } else {
          this.installText = resData.message
        }

        if (resData.speed == 100 || this.errorType == 3) {
          clearInterval(this.timer)
        }
        let _this = this
        if (resData.speed == 100) {
          setTimeout(() => {
            _this.$emit('updateState')
            _this.$emit('close')
          }, 1000)
        }
      })
    },

    /**
     * @description: Save edit update
     * @return {*} void
     */
    updateApp() {
      this.processData();
      this.isLoading = true;
      this.$api.app.updateContainerSetting(this.id, this.initData).then((res) => {
        if (res.data.success == 200) {
          this.isLoading = false;
          this.$emit('updateState')
        } else {
          this.$buefy.toast.open({
            message: res.data.message,
            type: 'is-warning'
          })
        }
        this.$emit('close')
      })
    },

    /**
     * @description: Show import panel
     * @return {*} void
     */
    showImportPanel() {
      this.$buefy.modal.open({
        parent: this,
        component: ImportPanel,
        hasModalCard: true,
        customClass: '',
        trapFocus: true,
        canCancel: ['escape'],
        scroll: "keep",
        animation: "zoom-out",
        events: {
          'update': (e) => {
            this.initData = e
            this.$buefy.dialog.alert({
              title: 'Attention',
              message: 'AutoFill only helps you to complete most of the configuration. Some of the configuration information still needs to be confirmed by you.',
              type: 'is-dark'
            })
          }
        },
        props: {
          initData: this.initData,
          netWorks: this.networks
        }
      })
    },

    /**
     * @description: Get remote synchronization information
     * @param {*} function
     * @return {*} void
     */
    getAsyncData: debounce(function (name) {
      if (!name.length) {
        this.data = []
        return
      }
      this.isFetching = true
      axios.get(`https://hub.docker.com/api/content/v1/products/search?source=community&q=${name}&page=1&page_size=4`)
        .then(({ data }) => {
          this.data = []
          data.summaries.forEach((item) => this.data.push(item.name))
        })
        .catch((error) => {
          this.data = []
          throw error
        })
        .finally(() => {
          this.isFetching = false
        })
    }, 500)

  },
  destroyed() {
    clearInterval(this.timer)
  },
}
</script>
