<!--
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-19 09:23:01
 * @Description: Top Search bar
 * @FilePath: \CasaOS-UI\src\components\SearchBar.vue
-->

<template>
  <b-field position="is-centered " class="search-bar has-text-white">
    <b-input placeholder="Google Search..." v-model="keyText" icon="magnify" icon-right="magnify" icon-right-clickable @icon-right-click="gotoSearch" @keyup.enter.native="gotoSearch" size="is-large" :class="['ovh',isFocus?'fo':'']" expanded @focus="onFocus" @blur="onBlur">
    </b-input>
  </b-field>
</template>

<script>
export default {
  name: "search-bar",
  data() {
    return {
      isFocus: false,
      keyText: ""
    }
  },
  methods: {
    /**
     * @description: Handle Focus event
     * @return {*} void
     */
    onFocus() {
      this.isFocus = true;
    },

    /**
     * @description: Handle Blur event
     * @return {*} void
     */
    onBlur() {
      if (this.keyText == "")
        this.isFocus = false;
    },

    /**
     * @description: Pop up a new window and jump to google search
     * @return {*} void
     */
    gotoSearch() {
      window.open("https://www.google.com/search?q=" + this.keyText, '_blank')
    }
  },
}
</script>

<style lang="scss">
.search-bar {
  input {
    transition: all 0.2s;
    appearance: none;
    background: rgba(123, 123, 123, 0.16);
    backdrop-filter: blur(0.875rem);
    border-radius: 8px;
    border: none;
    outline: none;
    font-size: 1.5rem;
    color: white;
    &:focus {
      border: none;
      box-shadow: none;
    }
    &::placeholder {
      color: white;
    }
  }
  .ovh {
    overflow: hidden;
    .icon.is-left {
      transition: all 0.2s;
      left: 0;
    }
    .icon.is-right {
      transition: all 0.2s;
      right: -3rem !important;
      color: white !important;
    }
    input {
      padding-left: 2.5em !important;
      padding-right: 1em !important;
    }
  }
  .fo {
    .icon.is-left {
      left: -3rem !important;
    }
    .icon.is-right {
      transition: all 0.2s;
      right: 0 !important;
    }
    input {
      padding-right: 2.5em !important;
      padding-left: 1em !important;
    }
  }
}
</style>