<template>
  <div class="widget has-text-white clock">
    <div class="time">{{timeText}}</div>
    <div class="data">{{dateText}}</div>
  </div>
</template>

<script>
import moment from 'moment'
export default {
  data() {
    return {
      timer: 0,
      timeText:"",
      dateText:""
    }
  },
  mounted() {
    if (this.timer) {
      clearInterval(this.timer)
    }
    this.updateClock()
    this.timer = setInterval(() => {
      this.updateClock()
    }, 1000)
  },
  methods: {
    updateClock() {
        
        this.timeText = moment().format('LT');
        this.dateText = moment().format('dddd, MMMM Do')
    }
  },
}
</script>

<style lang="scss">
.clock {
  font-family: Roboto;
  font-style: normal;
  font-weight: 300;
  text-align: left;
  .time {
    font-size: 2.75rem;
    line-height: 4.25rem;
    opacity: 0.9;
  }
  .data {
    line-height: 1.5rem;
    opacity: 0.9;
  }
}
</style>