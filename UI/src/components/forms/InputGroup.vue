<template>
  <div class="mb-5">
    <div class="field is-flex is-align-items-center mb-2">
      <label class="label mb-0 flex1">{{label}}</label>
      <b-button icon-left="plus" type="is-dark" size="is-small" rounded @click="addItem">Add</b-button>
    </div>
    <div class="is-flex is-align-items-center mb-5 info" v-if="vdata.length == 0">
      <b-icon pack="fas" icon="info-circle" size="is-small" class="mr-2 "></b-icon>
      <span>
        {{message}}
      </span>

    </div>

    <div class="port-item" v-for="(item,index) in vdata" :key="'port'+index">
      <b-icon pack="fas" icon="times" size="is-small" class="is-clickable" @click.native="removeItem(index)"></b-icon>
      <template v-if="index < 1">
        <b-field grouped>
          <b-field :label="name1" expanded>
            <b-input :placeholder="name1" v-model="item.container" expanded @input="handleInput"></b-input>
          </b-field>
          <b-field :label="name2" expanded>
            <b-input :placeholder="name2" v-model="item.host" expanded @input="handleInput"></b-input>
          </b-field>

        </b-field>
      </template>
      <template v-else>

        <b-field grouped>
            <b-input :placeholder="name1" v-model="item.container" expanded @input="handleInput"></b-input>
            <b-input :placeholder="name2" v-model="item.host" expanded @input="handleInput"></b-input>
        </b-field>

      </template>
    </div>

  </div>
</template>

<script>
export default {
  name:'input-group',
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
    label: String,
    message: String,
    name1: {
      type: String,
      default: "Container"
    },
    name2: {
      type: String,
      default: "Host"
    },

  },
  created() {
    //this.items = this.vdata;
  },
  watch: {

  },
  mounted() {
    //this.addItem()
  },
  methods: {
    addItem() {
      let itemObj = {
        container: "",
        host: ""
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

