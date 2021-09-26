<template>
  <div class="mb-5">
    <div class="field is-flex is-align-items-center mb-2">
      <label class="label mb-0 flex1">Ports</label>
      <b-button icon-left="plus" type="is-dark" size="is-small" rounded @click="addItem">Add</b-button>
    </div>
    <div class="is-flex is-align-items-center mb-5 info" v-if="vdata.length == 0">
      <b-icon pack="fas" icon="info-circle" size="is-small" class="mr-2 "></b-icon>
      <span>
        No App Ports now, Click “+” to add one.
      </span>

    </div>
    <div class="port-item" v-for="(item,index) in vdata" :key="'port'+index">
      <b-icon pack="fas" icon="times" size="is-small" class="is-clickable" @click.native="removeItem(index)"></b-icon>
      <template v-if="index < 1">
        <b-field grouped>
          <b-field label="Container" expanded>
            <b-input placeholder="Container" type="number" v-model="item.container" expanded @input="handleInput"></b-input>
          </b-field>
          <b-field label="Host" expanded>
            <b-input placeholder="Host" type="number" v-model="item.host" expanded @input="handleInput" v-if="showHostPost"></b-input>
          </b-field>
          <b-field label="Protocol" expanded>
            <b-select placeholder="Protocol" v-model="item.protocol" expanded @input="handleInput">
              <option value="tcp">TCP</option>
              <option value="udp">UDP</option>
              <option value="both">TCP + UDP</option>
            </b-select>
          </b-field>
        </b-field>
      </template>
      <template v-else>

        <b-field grouped>
          <b-input placeholder="Container" type="number" v-model="item.container" expanded @input="handleInput"></b-input>
          <b-input placeholder="Host" type="number" v-model="item.host" expanded @input="handleInput" v-if="showHostPost"></b-input>
          <b-select placeholder="Protocol" v-model="item.protocol" expanded @input="handleInput">
            <option value="tcp">TCP</option>
            <option value="udp">UDP</option>
            <option value="both">TCP + UDP</option>
          </b-select>
        </b-field>

      </template>
    </div>

  </div>
</template>

<script>
 
export default {
  name: 'ports',
  data() {
    return {
      isLoading: false,
      items: [],
      min: 0
    }
  },
  model: {
    prop: 'vdata',
    event: 'change'
  },
  props: {
    vdata: Array,
    showHostPost: Boolean
  },
  created() {
    //this.items = this.vdata;
  },
  mounted() {
    if (this.vdata.length == 0) {
      //this.addItem()
    }
  },
  methods: {
    addItem() {
      let itemObj = {
        container: "",
        host: "",
        protocol: "tcp"
      }
      this.vdata.push(itemObj)
    },

    removeItem(index) {
      this.vdata.splice(index, 1)
      this.filterArray()
    },
    handleInput() {
      this.filterArray()
    },
    filterArray() {
      // let newArray = this.items.filter(item => {
      //   if (item.container != "" && item.host != "") {
      //     return true
      //   } else {
      //     return false
      //   }
      // })

      this.$emit('change', this.vdata)


    }
  },
}
</script>

<style lang="scss">
.info {
  font-size: 0.875rem;
  color: #5a5a5a;
}
.port-item {
  position: relative;
  .icon {
    position: absolute;
    right: -1.5rem;
    bottom: 0.825rem;
  }
  &:not(:last-child) {
    margin-bottom: 0.5rem;
  }
  .field.is-expanded {
    .label {
      text-align: center;
      font-weight: normal;
    }
  }
}
</style>